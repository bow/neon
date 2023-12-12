// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/bow/lens/internal/database"
	"github.com/bow/lens/internal/tui"
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

			initPath, err := readerInitPath()
			if err != nil {
				return err
			}

			dbPath, err := resolveDBPath(v.GetString(dbPathKey))
			if err != nil {
				return err
			}
			fs, err := database.NewSQLite(dbPath)
			if err != nil {
				return err
			}

			app := tui.NewReader(cmd.Context(), fs).
				WithInitPath(initPath)

			return app.Show()
		},
	}

	pflags := command.PersistentFlags()

	pflags.StringP(dbPathKey, "d", defaultDBPath, "data store location")

	if err := v.BindPFlags(pflags); err != nil {
		panic(err)
	}

	return &command
}

func readerInitPath() (string, error) {
	sd, err := stateDir()
	if err != nil {
		return "", err
	}
	_, err = os.Stat(sd)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		if err := os.MkdirAll(sd, os.ModeDir|0o700); err != nil {
			return "", err
		}
	}
	return filepath.Join(sd, "reader.initialized"), nil
}
