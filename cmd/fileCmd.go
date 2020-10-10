package cmd

import (
	"github.com/spf13/cobra"
)

func NewFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "file <path>",
		Short: "",
		Long:  ``,
		Run: func(command *cobra.Command, args []string) {
			handleFileCommand(args)
		},
	}
}

func handleFileCommand(args []string) {
	fileScanner := FileScanner{Logging: NewLogging()}
	fileScanner.list(args[0])
}

func init() {
	fileCmd:= NewFileCmd()
	rootCmd.AddCommand(fileCmd)
}
