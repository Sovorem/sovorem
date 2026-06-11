package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&forceSubmit, "submit", "s", false, "shortcut flag՝ run-ից հետո միանգամից submit անելու համար")
	runCmd.Flags().BoolVar(&debugSubmission, "debug", false, "log անել submission-ի request/response debug output-ը")
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:    "run UUID",
	Args:   cobra.MatchAll(cobra.RangeArgs(1, 10)),
	Short:  "Run անել դասը՝ առանց submit անելու",
	PreRun: compose(requireUpdated, requireAuth),
	RunE:   submissionHandler,
}
