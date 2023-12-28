// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bow/neon/internal/entity"
)

func newFeedShowEntryCommand() *cobra.Command {
	var name = "show-entry"

	command := cobra.Command{
		Use:                   fmt.Sprintf("%s ENTRY-ID", name),
		Aliases:               []string{"show-e", "se"},
		Short:                 "Show a feed entry",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return fmt.Errorf("entry ID value not specified")
			} else if len(args) > 1 {
				return fmt.Errorf("too many arguments")
			}

			entryID, err := entity.ToFeedID(args[0])
			if err != nil {
				return err
			}

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			entry, err := db.GetEntry(cmd.Context(), entryID)
			if err != nil {
				return err
			}

			if content := entry.Content; content != nil {
				fmt.Printf("%s\n", *content)
			}

			return nil
		},
	}

	return &command
}
