package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
)

func (app *Application) NewScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "",
		Long:  ``,
		Run: func(command *cobra.Command, args []string) {
			app.handleScanCommand()
		},
	}
}

func (app *Application) handleScanCommand() {
	scanner := Scanner{PClient: app.PClient, Logging: app.Logging}
	scanResult, _ := scanner.scan("https://example.com", 10, 2)

	scanResultJson, err := json.Marshal(scanResult)

	if err != nil {
		app.Logging.errorLog.Printf("%v", err)
	}

	app.Logging.stdout.Printf(string(scanResultJson))
}

func init() {
	logging := NewLogging()
	pclient := NewPClient(*logging)

	app := NewApplication(logging, pclient)
	scanCmd := app.NewScanCmd()
	rootCmd.AddCommand(scanCmd)
}
