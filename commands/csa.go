package commands

import (
	"os"

	cfg "github.com/scalefast/talos/config"
	logger "github.com/scalefast/talos/logger"
	"github.com/scalefast/talos/tools/csa"
	"github.com/spf13/cobra"
)

// cmdCSA is a global variable that contains the information needed for cobra to
// add the command to the command-line parser.
var cmdCSA = &cobra.Command{
	Use:   "csa [string to echo]",
	Short: "Container scanner",
	Long:  `This application is used to scan a container looking for vulnerabilities`,
	Run: func(cmd *cobra.Command, args []string) {

		l := logger.NewLogger("csa")

		// Parse the config file specified by parameters. // Sergio's comment
		c, err := cfg.ParseConfig(filename)
		if err != nil {
			//TODO: Replace with log call
			l.ECustom("Config parsed wrongly")
			os.Exit(1)
		}
		// Send the control, (along with the config, and logger) // Sergio's comment
		// To the analyzer // Sergio's comment
		csa.Analyze(c, l, user, network, output)
	},
}
