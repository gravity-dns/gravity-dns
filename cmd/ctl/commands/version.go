package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gravity server.",
	Long:  "All software has versions. This is Gravity server's.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gravity DNS Controller v0.01 -- HEAD")
	},
}
