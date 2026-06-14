package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	api "github.com/Sovorem/sovorem-cli/client"
	"github.com/Sovorem/sovorem-cli/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	apiURLOverride string
)

var rootCmd = &cobra.Command{
	Use:   "sovorem",
	Short: "Sovorem.am-ի պաշտոնական CLI",
	Long: `Sovorem.am-ի պաշտոնական CLI-ը։ Էս ծրագիրը նախատեսված ա
որպես կայքի օգնական app (ոչ թե փոխարինող)։`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(currentVersion string) error {
	rootCmd.Version = currentVersion
	info := version.FetchUpdateInfo(rootCmd.Version)
	defer info.PromptUpdateIfAvailable()
	ctx := version.WithContext(context.Background(), &info)
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config ֆայլ (default-ը $HOME/.sovorem.yaml կամ $XDG_CONFIG_HOME/sovorem/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURLOverride, "api-url", "", "Sovorem-ի base URL-ը override անել միայն տվյալ command-ի համար, օր․ http://localhost:3000 (չի պահվում config-ում)")
}

func readViperConfig(paths []string) error {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			viper.SetConfigFile(p)
			break
		}
	}
	return viper.ReadInConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// The API/login base URL is environment config, not user data: it lives in
	// the binary (api.DefaultBaseURL) and is overridable via --api-url /
	// SOVOREM_API_URL. It is intentionally NOT a viper default, so it is never
	// written into the user's config file where a stale value could shadow it.
	viper.SetDefault("access_token", "")
	viper.SetDefault("refresh_token", "")
	viper.SetDefault("last_refresh", 0)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(filepath.Clean(cfgFile))
		cobra.CheckErr(viper.ReadInConfig())
	} else {
		// find home dir
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// collect paths where existing config files may be located
		var configPaths []string

		// first check XDG_CONFIG_HOME if set
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		var xdgEnvPath string
		if xdgConfigHome != "" {
			xdgEnvPath = filepath.Join(xdgConfigHome, "sovorem", "config.yaml")
			configPaths = append(configPaths, xdgEnvPath)
		}

		// then check legacy hard-coded "XDG" path, then home dotfile
		xdgLegacyPath := filepath.Join(home, ".config", "sovorem", "config.yaml")
		homeDotfilePath := filepath.Join(home, ".sovorem.yaml")

		configPaths = append(configPaths, xdgLegacyPath)
		configPaths = append(configPaths, homeDotfilePath)

		if err := readViperConfig(configPaths); err != nil {
			// no existing config found; try to create a new one
			// respect XDG_CONFIG_HOME if set, otherwise use dotfile in home dir
			var newConfigPath string
			if xdgEnvPath != "" {
				newConfigPath = xdgEnvPath
				cobra.CheckErr(os.MkdirAll(filepath.Dir(newConfigPath), 0o755))
			} else {
				newConfigPath = homeDotfilePath
			}

			cobra.CheckErr(viper.SafeWriteConfigAs(newConfigPath))
			viper.SetConfigFile(newConfigPath)
			cobra.CheckErr(viper.ReadInConfig())
		}
	}

	// Drop endpoint keys that older CLI versions persisted into the config file,
	// so a saved value can't shadow the binary default after an upgrade.
	migrateConfig()

	viper.SetEnvPrefix("sovorem")
	viper.AutomaticEnv() // read in environment variables that match

	// --api-url overrides both the API endpoint and the browser login URL for
	// this invocation only (handy for pointing the CLI at a local or staging
	// server). It is never written to the config file.
	cobra.CheckErr(api.SetBaseOverride(apiURLOverride))
}

// Chain multiple commands together.
func compose(commands ...func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, command := range commands {
			command(cmd, args)
		}
	}
}

// Call this function at the beginning of a command handler
// if you want to require the user to update their CLI first.
func requireUpdated(cmd *cobra.Command, args []string) {
	info := version.FromContext(cmd.Context())
	if info == nil {
		if !promptToContinue(
			"ԶԳՈՒՇԱՑՈՒՄ. Հնարավոր չի ստանալ version-ի տվյալները",
			"Հնարավոր չեղավ ստուգել՝ արդյոք քո sovorem CLI-ը update արած ա։",
			"Շարունակե՞նք ամեն դեպքում",
		) {
			os.Exit(1)
		}
		return
	}
	if info.FailedToFetch != nil {
		if !promptToContinue(
			"ԶԳՈՒՇԱՑՈՒՄ. Հնարավոր չի ստանալ version-ի տվյալները",
			fmt.Sprintf("Հնարավոր չեղավ ստուգել՝ արդյոք քո sovorem CLI-ը update արած ա. %s", info.FailedToFetch.Error()),
			"Շարունակե՞նք ամեն դեպքում",
		) {
			os.Exit(1)
		}
		return
	}
	if info.IsUpdateRequired {
		info.PromptUpdateIfAvailable()
		os.Exit(1)
	}
}

func promptToContinue(title string, message string, prompt string) bool {
	fmt.Fprintln(os.Stderr, title)
	fmt.Fprintln(os.Stderr, message)
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// Call this function at the beginning of a command handler
// if you need to make authenticated requests. This will
// automatically refresh the tokens, if necessary, and prompt
// the user to re-login if anything goes wrong.
func requireAuth(cmd *cobra.Command, args []string) {
	promptLoginAndExitIf := func(condition bool) {
		if condition {
			fmt.Fprintln(os.Stderr, "Էս command-ը run անելու համար պետք ա login եղած լինես։")
			fmt.Fprintln(os.Stderr, "Արի սկզբից 'sovorem login' run անենք։")
			os.Exit(1)
		}
	}

	accessToken := viper.GetString("access_token")
	promptLoginAndExitIf(accessToken == "")

	// We only refresh if our token is getting stale.
	lastRefresh := viper.GetInt64("last_refresh")
	if time.Now().Add(-time.Minute*55).Unix() <= lastRefresh {
		return
	}

	creds, err := api.FetchAccessToken()
	promptLoginAndExitIf(err != nil)
	if creds.AccessToken == "" || creds.RefreshToken == "" {
		promptLoginAndExitIf(err != nil)
	}

	viper.Set("access_token", creds.AccessToken)
	viper.Set("refresh_token", creds.RefreshToken)
	viper.Set("last_refresh", time.Now().Unix())

	err = viper.WriteConfig()
	promptLoginAndExitIf(err != nil)
}
