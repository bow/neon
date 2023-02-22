package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newFeedCmd() *cobra.Command {
	var (
		name      = "feed"
		feedViper = newViper(name)
	)

	feed := cobra.Command{
		Use:   name,
		Short: "View or modify feeds",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			dbPath, err := resolveDBPath(feedViper)
			if err != nil {
				return err
			}
			ctx := context.WithValue(cmd.Context(), ctxKey(dbNameKey), dbPath)
			cmd.SetContext(ctx)

			return nil
		},
	}

	pflags := feed.PersistentFlags()

	pflags.StringP(dbNameKey, "d", defaultDBName, "data store location")
	_ = feedViper.BindPFlag(dbNameKey, pflags.Lookup(dbNameKey))

	feed.AddCommand(newFeedListCmd())

	return &feed
}
