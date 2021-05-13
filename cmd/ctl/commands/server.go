package commands

import (
	"github.com/gravity-dns/gravity-dns/dns"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startServerCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the Gravity DNS server.",
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Gravity DNS server.",
	Run: func(cmd *cobra.Command, args []string) {
		dns.StartServer()
	},
}
