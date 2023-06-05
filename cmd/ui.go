// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUICmd() *cobra.Command {
	var (
		name    = "ui"
		uiViper = newViper(name)
	)

	uiCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Open a TUI feed reader",
		RunE: func(cmd *cobra.Command, args []string) error {

			dbPath, err := resolveDBPath(uiViper.GetString(dbPathKey))
			if err != nil {
				return err
			}
			dataStorePathToCmdCtx(cmd, dbPath)

			return fmt.Errorf("not fully implemented")
		},
	}

	pflags := uiCmd.PersistentFlags()

	pflags.StringP(dbPathKey, "d", defaultDBPath, "data store location")

	if err := uiViper.BindPFlags(pflags); err != nil {
		panic(err)
	}

	return &uiCmd
}
