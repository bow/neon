package cmd

import (
	"github.com/adrg/xdg"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bow/courier/internal"
)

const (
	addrKey   = "addr"
	dbNameKey = "db"
)

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		store, err := internal.NewFeedsDB(viper.GetString(dbNameKey))
		if err != nil {
			return err
		}

		server, err := internal.NewServerBuilder().
			Address(viper.GetString(addrKey)).
			Store(store).
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

	// TODO: Handle error.
	defaultDBName, _ := xdg.DataFile("courier/courier.db")
	flags.StringP(dbNameKey, "d", defaultDBName, "data store location")
	_ = viper.BindPFlag(dbNameKey, flags.Lookup(dbNameKey))
}
