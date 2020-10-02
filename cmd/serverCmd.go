package cmd

import (
	"github.com/spf13/cobra"
)

func NewServerScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "",
		Long:  ``,
		Run: func(command *cobra.Command, args []string) {
			handleServerCommand()
		},
	}
}

func handleServerCommand() {
	server := Server{Logging: NewLogging()}
	server.start()
}

func init() {
	serverCmd := NewServerScanCmd()
	rootCmd.AddCommand(serverCmd)
}
