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
		log.Info("starting server")
		srv := server.New()
		err := srv.Serve()
		if err != nil {
			log.Error(err)
		}
	},
}
