package cmd

import (
	"github.com/bow/courier/server"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const addrKey = "addr"

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := viper.GetString(addrKey)

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

	flags := serveCmd.Flags()

	flags.StringP(addrKey, "a", ":50051", "listening address")
	_ = viper.BindPFlag(addrKey, flags.Lookup(addrKey))
}
