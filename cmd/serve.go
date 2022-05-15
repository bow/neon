package cmd

import (
	"github.com/bow/courier/server"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Expose as proper flag.
		addr := ":50051"
		logger := zlog.Logger.With().Logger()

		builder := server.NewBuilder().
			Address(addr).
			Logger(logger)

		server, err := builder.Build()
		if err != nil {
			return err
		}

		return server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(&serveCmd)
}
