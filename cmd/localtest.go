package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sovorem/sovorem/checks"
	api "github.com/sovorem/sovorem/client"
	"github.com/sovorem/sovorem/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

func init() {
	rootCmd.AddCommand(localTestCmd)
}

var localTestCmd = &cobra.Command{
	Use:    "local-test PATH",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	RunE:   localTestHandler,
}

func localTestHandler(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	data, err := readLocalCLIData(args[0])
	if err != nil {
		return err
	}
	if err := validateAllowedOS(data); err != nil {
		return err
	}

	overrideBaseURL := viper.GetString("override_base_url")
	if overrideBaseURL != "" {
		fmt.Printf("Օգտագործվում ա override արած base_url-ը. %v\n", overrideBaseURL)
		fmt.Printf("Default-ին կարող ես վերադառնալ `sovorem config base_url --reset` run անելով\n\n")
	}

	ch := make(chan tea.Msg, 1)
	finalise := render.StartRenderer(data, true, ch)

	cliResults := checks.CLIChecks(data, overrideBaseURL, ch)
	submissionEvent := checks.LocalSubmissionEvent(data, cliResults)
	checks.ApplySubmissionResults(data, submissionEvent.StructuredErrCLI, ch)
	finalise(submissionEvent)

	if submissionEvent.ResultSlug != api.VerificationResultSlugSuccess {
		return localTestFailureError(submissionEvent.StructuredErrCLI)
	}

	return nil
}

func localTestFailureError(failure *api.StructuredErrCLI) error {
	if failure == nil {
		return errors.New("լոկալ test-երը չանցան")
	}
	return fmt.Errorf(
		"լոկալ test-երը չանցան. step %d, test %d\n%s",
		failure.FailedStepIndex+1,
		failure.FailedTestIndex+1,
		failure.ErrorMessage,
	)
}

func readLocalCLIData(path string) (api.CLIData, error) {
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return api.CLIData{}, err
	}
	if info.IsDir() {
		cleanPath = filepath.Join(cleanPath, "cli.yaml")
	}

	bytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return api.CLIData{}, err
	}

	var data api.CLIData
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return api.CLIData{}, err
	}
	if len(data.Steps) == 0 {
		return api.CLIData{}, errors.New("test manifest-ը պետք ա առնվազն մեկ step ներառի")
	}

	return data, nil
}

func validateAllowedOS(data api.CLIData) error {
	if len(data.AllowedOperatingSystems) == 0 {
		return nil
	}

	if slices.Contains(data.AllowedOperatingSystems, runtime.GOOS) {
		return nil
	}

	return fmt.Errorf(
		"դասը չի աջակցվում քո օպերացիոն համակարգի կողմից (%s)\nնորից փորձիր սրանցից մեկով. %v",
		runtime.GOOS,
		data.AllowedOperatingSystems,
	)
}
