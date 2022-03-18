package commands

import (
	"os"

	"github.com/spf13/cobra"
	cfg "github.com/scalefast/talos/config"
	logger "github.com/scalefast/talos/logger"
	"github.com/scalefast/talos/tools/dast"
)

var cmdDAST = &cobra.Command{
	Use:   "dast (--config config_file.[yaml|json])",
	Short: "Analyze an application dynamically",
	Long:  `This application is used to execute a dynamic security vulnerabilities analyzer.`,
	Run: func(cmd *cobra.Command, args []string) {

		l := logger.NewLogger("dast")

		c, err := cfg.ParseConfig(filename)
		if err != nil {
			l.EInvalidArg("Config file")
			os.Exit(1)
		}
		// Send the control, (along with the config, and logger)
		// To the analyzer
		dast.Analyze(c, l)
	},
}
