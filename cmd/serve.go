package cmd

import (
	"fmt"

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

var (
	relDBName     = "courier/courier.db"
	defaultDBName = fmt.Sprintf("$XDG_DATA_HOME/%s", relDBName)
)

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		dbPath, err := resolveDBPath()
		if err != nil {
			return err
		}

		store, err := internal.NewFeedsDB(dbPath)
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

	flags.StringP(dbNameKey, "d", defaultDBName, "data store location")
	_ = viper.BindPFlag(dbNameKey, flags.Lookup(dbNameKey))
}

func resolveDBPath() (dbPath string, err error) {
	dbPath = viper.GetString(dbNameKey)
	if dbPath == defaultDBName {
		dbPath, err = xdg.DataFile(relDBName)
		if err != nil {
			return "", err
		}
	}
	return dbPath, nil
}
