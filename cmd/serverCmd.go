package cmd

import (
	"github.com/spf13/cobra"
)

func NewServerScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server <port>",
		Short: "",
		Long:  ``,
		Run: func(command *cobra.Command, args []string) {
			handleServerCommand(args)
		},
	}
}

func handleServerCommand(args []string) {
	server := Server{Logging: NewLogging()}
	server.start(args[0])
}

func init() {
	serverCmd := NewServerScanCmd()
	rootCmd.AddCommand(serverCmd)
}
