package cmd

import (
	"fmt"

	api "github.com/Sovorem/sovorem/client"
	"github.com/Sovorem/sovorem/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Տեսնել login-ի և CLI-ի version-ի status-ը",
	Long:  "Ցույց տալ՝ արդյոք login ես եղել և արդյոք CLI-ը update արած ա",
	Run: func(cmd *cobra.Command, args []string) {
		checkAuthStatus()
		fmt.Println() // Blank line for readability
		checkVersionStatus(cmd)
	},
}

func checkAuthStatus() {
	refreshToken := viper.GetString("refresh_token")
	if refreshToken == "" {
		fmt.Println("Login եղած չես")
		fmt.Println("Run արա 'sovorem login'՝ login լինելու համար")
		return
	}

	// Verify token is still valid by attempting to refresh
	_, err := api.FetchAccessToken()
	if err != nil {
		fmt.Println("Authentication-ի ժամկետն անցել ա")
		fmt.Println("Run արա 'sovorem login'՝ նորից login լինելու համար")
		return
	}

	user, err := api.FetchCurrentUser()
	if err == nil && user != nil && user.Handle != "" {
		fmt.Printf("Login ես եղել որպես @%s\n", user.Handle)
		return
	}

	fmt.Println("Login եղած ես")
}

func checkVersionStatus(cmd *cobra.Command) {
	info := version.FromContext(cmd.Context())
	if info == nil || info.FailedToFetch != nil {
		fmt.Println("Հնարավոր չեղավ ստուգել version-ի status-ը")
		if info != nil && info.FailedToFetch != nil {
			fmt.Printf("Error: %s\n", info.FailedToFetch.Error())
		}
		return
	}

	if info.IsOutdated {
		fmt.Printf("CLI-ը հնացել ա. %s → հասանելի ա %s\n", info.CurrentVersion, info.LatestVersion)
		fmt.Println("Run արա 'sovorem upgrade'՝ update անելու համար")
	} else {
		fmt.Printf("CLI-ը update արած ա (%s)\n", info.CurrentVersion)
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
