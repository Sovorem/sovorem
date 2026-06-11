package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sovorem/sovorem/checks"
	api "github.com/sovorem/sovorem/client"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sovorem/sovorem/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	forceSubmit     bool
	debugSubmission bool
)

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVar(&debugSubmission, "debug", false, "log անել submission-ի request/response debug output-ը")
}

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:    "submit UUID",
	Args:   cobra.MatchAll(cobra.RangeArgs(1, 10)),
	Short:  "Submit անել դասը",
	PreRun: compose(requireUpdated, requireAuth),
	RunE:   submissionHandler,
}

func submissionHandler(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	isSubmit := cmd.Name() == "submit" || forceSubmit
	lessonUUID := args[0]

	lesson, err := api.FetchLesson(lessonUUID)
	if err != nil {
		return err
	}
	if lesson.Lesson.Type != "type_cli" {
		return errors.New("հնարավոր չի run անել դասը. չաջակցվող դասի տիպ")
	}
	if lesson.Lesson.LessonDataCLI == nil {
		return errors.New("հնարավոր չի run անել դասը. դասի տվյալները բացակայում են")
	}

	data := lesson.Lesson.LessonDataCLI.CLIData

	isAllowedOS := false
	for _, system := range data.AllowedOperatingSystems {
		if system == runtime.GOOS {
			isAllowedOS = true
		}
	}

	if !isAllowedOS {
		return fmt.Errorf("դասը չի աջակցվում քո օպերացիոն համակարգի կողմից (%s), նորից փորձիր սրանցից մեկով. %v", runtime.GOOS, data.AllowedOperatingSystems)
	}

	overrideBaseURL := viper.GetString("override_base_url")
	if overrideBaseURL != "" {
		fmt.Printf("Օգտագործվում ա override արած base_url-ը. %v\n", overrideBaseURL)
		fmt.Printf("Default-ին կարող ես վերադառնալ `sovorem config base_url --reset` run անելով\n\n")
	}

	ch := make(chan tea.Msg, 1)
	// StartRenderer and returns immediately, finalise function blocks the execution until the renderer is closed.
	finalise := render.StartRenderer(data, isSubmit, ch)

	cliResults := checks.CLIChecks(data, overrideBaseURL, ch)

	if isSubmit {
		submissionEvent, debugData, err := api.SubmitCLILesson(lessonUUID, cliResults, debugSubmission)
		if debugSubmission {
			var debugPath string
			var debugWriteErr error
			defer func() {
				reportDebugFileWrite(debugPath, debugWriteErr)
			}()
			debugPath, debugWriteErr = writeSubmissionDebugFile(lessonUUID, debugData)
		}
		if err != nil {
			return err
		}
		checks.ApplySubmissionResults(data, submissionEvent.StructuredErrCLI, ch)
		finalise(submissionEvent)
	} else {
		finalise(api.LessonSubmissionEvent{})
	}
	return nil
}

func reportDebugFileWrite(path string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "զգուշացում. չհաջողվեց գրել submission-ի debug output-ը. %v\n", err)
		return
	}
	fmt.Fprintf(os.Stderr, "Submission-ի debug output-ը գրվեց %s ֆայլում\n", path)
}

func writeSubmissionDebugFile(lessonUUID string, data api.SubmissionDebugData) (string, error) {
	now := time.Now()
	timestamp := now.Format("20060102-150405")
	filename := fmt.Sprintf("sovorem-submit-debug-%s-%s.txt", lessonUUID, timestamp)
	status := "unavailable"
	if data.ResponseStatusCode != 0 {
		status = fmt.Sprintf("%d", data.ResponseStatusCode)
	}

	contents := fmt.Sprintf(
		"sovorem submit debug\nTimestamp: %s\nLesson UUID: %s\nEndpoint: %s\n\n=== Request JSON ===\n%s\n\n=== Response ===\nStatus Code: %s\n%s\n",
		now.Format(time.RFC3339),
		lessonUUID,
		data.Endpoint,
		data.RequestBody,
		status,
		data.ResponseBody,
	)

	if err := os.WriteFile(filename, []byte(contents), 0o600); err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return filename, nil
	}

	return absPath, nil
}
