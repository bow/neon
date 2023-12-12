// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bow/lens/internal"
)

func newFeedListEntriesCommand() *cobra.Command {
	var (
		name = "list-entries"
		v    = newViper(name)
	)

	const (
		bookmarkedKey = "bookmarked"
	)

	command := cobra.Command{
		Use:                   fmt.Sprintf("%s FEED-ID", name),
		Aliases:               []string{"list-e", "le"},
		Short:                 "List feed entries",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return fmt.Errorf("feed ID value not specified")
			} else if len(args) > 1 {
				return fmt.Errorf("too many arguments")
			}

			var isBookmarked *bool = nil
			if value := v.GetBool(bookmarkedKey); value {
				isBookmarked = &value
			}

			feedID, err := internal.ToFeedID(args[0])
			if err != nil {
				return err
			}

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			entries, err := db.ListEntries(cmd.Context(), []internal.ID{feedID}, isBookmarked)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				fmt.Printf("%s\n", fmtListEntry(entry))
			}

			return nil
		},
	}

	flags := command.Flags()

	flags.BoolP(bookmarkedKey, "b", false, "list only bookmarked entries")

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	return &command
}

func fmtListEntry(entry *internal.Entry) string {
	var (
		sb  strings.Builder
		cat = func(format string, a ...any) { fmt.Fprintf(&sb, format, a...) }
	)

	kv := []*struct {
		k, v string
	}{
		{"EntryID", fmt.Sprintf("%d", entry.ID)},
		{"URL", derefOrEmpty(entry.URL)},
		{"Pub", fmtOrEmpty(entry.Published)},
	}

	keyMaxLen := 0
	for _, line := range kv {
		keyLen := len(line.k)
		if keyLen > keyMaxLen {
			keyMaxLen = keyLen
		}
	}

	cat("\x1b[36m▶\x1b[0m \x1b[4m%s\x1b[0m\n", capText(entry.Title))
	for _, line := range kv {
		if line.v == "" {
			continue
		}
		cat("  %*s : %s\n", -1*keyMaxLen, line.k, capText(line.v))
	}

	return sb.String()
}
