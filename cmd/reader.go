// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/store"
	"github.com/bow/iris/internal/tui"
)

func newReaderCommand() *cobra.Command {
	var (
		name = "reader"
		v    = newViper(name)
	)

	command := cobra.Command{
		Use:     name,
		Aliases: append(makeAlias(name), []string{"ui", "tui"}...),
		Short:   "Open a feed reader",
		RunE: func(cmd *cobra.Command, args []string) error {

			dbPath, err := resolveDBPath(v.GetString(dbPathKey))
			if err != nil {
				return err
			}
			fs, err := store.NewSQLite(dbPath)
			if err != nil {
				return err
			}

			return tui.Show(fs)
		},
	}

	pflags := command.PersistentFlags()

	pflags.StringP(dbPathKey, "d", defaultDBPath, "data store location")

	if err := v.BindPFlags(pflags); err != nil {
		panic(err)
	}

	return &command
}
