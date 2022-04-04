package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"service_admin_contractor/application"
)

// Является serve командой, ответственной за запуск API
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts http server with configured api",
	Long:  `Starts http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := application.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
