package cmd

import (
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bow/courier/internal"
)

const addrKey = "addr"

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		server, err := internal.NewServerBuilder().
			Address(viper.GetString(addrKey)).
			Logger(zlog.Logger.With().Logger()).
			Build()

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
