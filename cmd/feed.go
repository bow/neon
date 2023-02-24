package cmd

import (
	"github.com/bow/iris/internal/store"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func newFeedCmd() *cobra.Command {
	var (
		name      = "feed"
		feedViper = newViper(name)
	)

	feedCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "View or modify feeds",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if zerolog.GlobalLevel() == zerolog.InfoLevel {
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			}

			dbPath, err := resolveDBPath(feedViper.GetString(dbPathKey))
			if err != nil {
				return err
			}
			inCmdContext(cmd, dbPathKey, dbPath)

			return nil
		},
	}

	pflags := feedCmd.PersistentFlags()

	pflags.StringP(dbPathKey, "d", defaultDBPath, "data store location")
	_ = feedViper.BindPFlag(dbPathKey, pflags.Lookup(dbPathKey))

	feedCmd.AddCommand(newFeedListCmd())
	feedCmd.AddCommand(newFeedExportCmd())

	return &feedCmd
}

func storeFromCtx(cmd *cobra.Command) (*store.SQLite, error) {
	dbPath, err := fromCmdContext[string](cmd, dbPathKey)
	if err != nil {
		return nil, err
	}
	str, err := store.NewSQLite(dbPath)
	if err != nil {
		return nil, err
	}
	return str, nil
}
