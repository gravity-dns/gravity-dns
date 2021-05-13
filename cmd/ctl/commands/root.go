package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gravity",
	Short: "Open source, edge DNS.",
	Long:  "Gravity DNS is an open source edge DNS server with a focus on privacy.",
}

func Run(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
