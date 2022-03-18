package commands

import (
	"github.com/spf13/cobra"
)

//Config file. Declare type
var filename string
var filetype string

// Gets called with the version we want to have,
// Now we want to print the version if it is required.
func Execute(version string) {
	cmdDAST.Flags().StringVarP(&filename, "config", "c", "config.yaml", "Config file location, relative to binary") // What handler handles file types in cobra?
	cmdCSA.Flags().StringVarP(&filename, "config", "c", "config.yaml", "Config file location, relative to binary")  // What handler handles file types in cobra?
	cmdAll.Flags().StringVarP(&filetype, "format", "f", "yaml", "Config file format")

	// Execute sets-up the run command, so that any security analysis tools can be
	// added after it.
	// It is the main program entrypoint, as all commands diverge from this point.
	var rootCmd = &cobra.Command{
		Use:   "talos",
		Short: "Perform automated security analysis on your application",
		Long: `This application performs automated security analysis on your application.
It's goal is to simplify security checks for developers, enabling them to
create robust, secure and quality code.
This program contains a series of security analysis tools, ranging from
     * Dynamic application security testing (DAST)
     * Container Security Analysis (CSA)

This tool has been developed by the security team at Scalefast.`,
		Version: version,
	}

	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdGen)
	cmdGen.AddCommand(cmdAll)
	cmdRun.AddCommand(cmdCSA)
	cmdRun.AddCommand(cmdDAST)
	// Execute the cobra command-line-parser
	rootCmd.Execute()

}
