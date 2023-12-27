// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bow/neon/internal/database"
)

func newFeedCommand() *cobra.Command {
	var (
		name = "feed"
		v    = newViper(name)
	)

	command := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "View or modify feeds",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			dbPath, err := resolveDBPath(v.GetString(dbPathKey))
			if err != nil {
				return err
			}
			dataStorePathToCmdCtx(cmd, dbPath)

			return nil
		},
	}

	pflags := command.PersistentFlags()

	pflags.StringP(dbPathKey, "d", defaultDBPath, "data store location")

	if err := v.BindPFlags(pflags); err != nil {
		panic(err)
	}

	command.AddCommand(newFeedAddCommand())
	command.AddCommand(newFeedExportCommand())
	command.AddCommand(newFeedImportCommand())
	command.AddCommand(newFeedListCommand())
	command.AddCommand(newFeedPullCommand())
	command.AddCommand(newFeedListEntriesCommand())
	command.AddCommand(newFeedShowEntryCommand())

	return &command
}

func dataStorePathToCmdCtx(cmd *cobra.Command, path string) {
	toCmdContext(cmd, dbPathKey, path)
}

func dbFromCmdCtx(cmd *cobra.Command) (*database.SQLite, error) {
	dbPath, err := fromCmdContext[string](cmd, dbPathKey)
	if err != nil {
		return nil, err
	}
	str, err := database.NewSQLite(dbPath)
	if err != nil {
		return nil, err
	}
	return str, nil
}
