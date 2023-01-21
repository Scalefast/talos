package commands

import (
	logger "github.com/scalefast/talos/logger"
	"github.com/spf13/cobra"
)

// Base command to RUN.
// Normally, we execute another command after, like dast
// But we will reserve this to include further support,
// like a global test.
var cmdRun = &cobra.Command{
	Use:   "run [tool to run]",
	Short: "Runs the specified security application",
	Run: func(cmd *cobra.Command, args []string) {
		l := logger.NewLogger("run")
		l.ICustom("Yet to be implemented!")
	},
}
