package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/Sovorem/sovorem/version"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Install անել CLI-ի ամենավերջին version-ը",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.FromContext(cmd.Context())
		if !info.IsOutdated {
			fmt.Println("Sovorem.am CLI-ը արդեն update արած ա։")
			return
		}

		fmt.Println("Update ենք անում Sovorem.am CLI-ը...")

		command := exec.Command("go", "install", "github.com/Sovorem/sovorem@latest")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()
		cobra.CheckErr(err)

		// Get new version info
		command = exec.Command("sovorem", "--version")
		versionBytes, err := command.Output()
		cobra.CheckErr(err)

		re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
		newVersion := re.FindString(string(versionBytes))
		if newVersion == "" {
			newVersion = "latest"
		}

		fmt.Printf("Հաջողությամբ update եղավ %s version-ին! 🎉\n", newVersion)
		os.Exit(0) // in case old version is still running
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
