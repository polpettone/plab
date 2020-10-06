package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
)

func (app *Application) NewScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan <url> <request count> <concurrency limit>",
		Short: "",
		Long:  ``,
		Run: func(command *cobra.Command, args []string) {
			app.handleScanCommand(command, args)
		},
	}
}

func (app *Application) handleScanCommand(cobraCommand *cobra.Command, args []string) {
	verbose, _ := cobraCommand.Flags().GetBool("verbose")
	infinite, _ := cobraCommand.Flags().GetBool("infinite")

	scanner := Scanner{PClient: app.PClient, Logging: app.Logging}

	if infinite {
		infiniteMode()
	} else {
		requestCount, _ := strconv.Atoi(args[1])
		concurrencyLimit, _ := strconv.Atoi(args[2])
		scanResult, _ := scanner.scan(args[0], requestCount, concurrencyLimit)
		if verbose {
			scanResultJsonVerbose, err := json.Marshal(scanResult)
			if err != nil {
				app.Logging.ErrorLog.Printf("%v", err)
			}
			app.Logging.Stdout.Printf(string(scanResultJsonVerbose))
		} else {
			scanResult.PClientResponses = nil
			scanResultJson, err := json.Marshal(scanResult)
			if err != nil {
				app.Logging.ErrorLog.Printf("%v", err)
			}
			app.Logging.Stdout.Printf(string(scanResultJson))
		}
	}
}

func infiniteMode() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()
	for {
		fmt.Println("do")
	}
}

func init() {
	logging := NewLogging()
	pclient := NewPClient(*logging)

	app := NewApplication(logging, pclient)
	scanCmd := app.NewScanCmd()

	scanCmd.Flags().BoolP(
		"verbose",
		"v",
		false,
		"Output includes all Responses",
	)

	scanCmd.Flags().BoolP(
		"infinite",
		"i",
		false,
		"run infinite mode, cancel with ctrl+c",
	)

	rootCmd.AddCommand(scanCmd)
}
