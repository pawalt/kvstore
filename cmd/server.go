package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/pawalt/kvstore/pkg/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start in server mode",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("expected a filepath but got none")
		}

		log.Info("starting server")
		srv, err := server.New(args[0])
		if err != nil {
			log.Fatal(err)
		}

		err = srv.Serve()
		if err != nil {
			log.Fatal(err)
		}
	},
}
