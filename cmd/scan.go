package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"strconv"
)

func (app *Application) NewScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan <url> <request count> <concurrency limit>",
		Short: "",
		Long:  ``,
		Run: func(command *cobra.Command, args []string) {
			app.handleScanCommand(args)
		},
	}
}

func (app *Application) handleScanCommand(args []string) {
	scanner := Scanner{PClient: app.PClient, Logging: app.Logging}

	requestCount, _  := strconv.Atoi(args[1])
	concurrencyLimit, _  := strconv.Atoi(args[2])

	scanResult, _ := scanner.scan(args[0], requestCount, concurrencyLimit)

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
