package cmd

import (
	"github.com/spf13/cobra"
)

func (app *Application) NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleRunCommand()
		},
	}
}

func (app *Application) handleRunCommand() {
	app.Logging.infoLog.Printf("run command")
	app.Logging.stdout.Printf("run command")
}


func init() {
	logging := NewLogging()
	pclient := NewPClient(*logging)
	app := NewApplication(logging, pclient)
	runCmd := app.NewRunCmd()
	rootCmd.AddCommand(runCmd)
}





