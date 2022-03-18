package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdGen = &cobra.Command{
	Use:   "gen [tool to generate]",
	Short: "generate the settings",
	Long:  `generates the settings required for talos to run`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Yet to be implemented")
	},
}
