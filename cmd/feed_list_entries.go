// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/store"
)

func newFeedListEntriesCommand() *cobra.Command {
	var name = "list-entries"

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

			feedID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			str, err := dataStoreFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			entries, err := str.ListEntries(cmd.Context(), store.ID(feedID))
			if err != nil {
				return err
			}
			for _, entry := range entries {
				fmt.Printf("%s\n", fmtListEntry(entry))
			}

			return nil
		},
	}

	return &command
}

func fmtListEntry(entry *store.Entry) string {
	var (
		sb  strings.Builder
		cat = func(format string, a ...any) { fmt.Fprintf(&sb, format, a...) }
	)

	var pubs = ""
	pub, err := deserializeTime(&entry.Published.String)
	if err != nil {
		pub = nil
	}
	if pub != nil {
		pubs = pub.Local().Format("2 January 2006 • 15:04 MST")
	}

	kv := []*struct {
		k, v string
	}{
		{"EntryID", fmt.Sprintf("%d", entry.ID)},
		{"URL", entry.URL.String},
		{"Pub", pubs},
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
