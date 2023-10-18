// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/store"
)

func newFeedViewEntryCmd() *cobra.Command {
	var name = "view-entry"

	viewCmd := cobra.Command{
		Use:                   fmt.Sprintf("%s ENTRY-ID", name),
		Aliases:               []string{"view-e", "ve"},
		Short:                 "View a feed entry",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return fmt.Errorf("entry ID value not specified")
			} else if len(args) > 1 {
				return fmt.Errorf("too many arguments")
			}

			entryID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			str, err := dataStoreFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			entry, err := str.ViewEntry(cmd.Context(), store.DBID(entryID))
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", entry.Content.String)

			return nil
		},
	}

	return &viewCmd
}
