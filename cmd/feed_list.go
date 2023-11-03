// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/store"
)

func newFeedListCommand() *cobra.Command {
	var name = "list"

	command := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "List feeds",

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := dataStoreFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			feeds, err := str.ListFeeds(cmd.Context())
			if err != nil {
				return err
			}
			for _, feed := range feeds {
				fmt.Printf("%s", fmtFeed(feed))
			}

			return nil
		},
	}

	return &command
}

func fmtFeed(feed *store.Feed) string {
	var (
		sb  strings.Builder
		cat = func(format string, a ...any) { fmt.Fprintf(&sb, format, a...) }
	)

	upds, _ := refmtTime(&feed.Updated.String)
	ps, _ := refmtTime(&feed.LastPulled)

	var nread, ntotal int
	for _, entry := range feed.Entries {
		if entry.IsRead {
			nread++
		}
		ntotal++
	}

	kv := []*struct {
		k, v string
	}{
		{"FeedID", fmt.Sprintf("%d", feed.ID)},
		{"Last pulled", ps},
		{"Updated", upds},
		{"Unread", fmt.Sprintf("%d/%d", ntotal-nread, ntotal)},
		{"URL", feed.SiteURL.String},
		{"Tags", fmt.Sprintf("#%s", strings.Join(feed.Tags, " #"))},
	}

	keyMaxLen := 0
	for _, line := range kv {
		keyLen := len(line.k)
		if keyLen > keyMaxLen {
			keyMaxLen = keyLen
		}
	}

	cat("\x1b[36m▶\x1b[0m \x1b[4m%s\x1b[0m\n", capText(feed.Title))
	for _, line := range kv {
		if line.v == "" {
			continue
		}
		cat("  %*s : %s\n", -1*keyMaxLen, line.k, capText(line.v))
	}
	cat("\n")

	return sb.String()
}

func refmtTime(raw *string) (rv string, err error) {
	deser, err := store.DeserializeTime(raw)
	if err != nil {
		deser = nil
	}
	if deser != nil {
		rv = deser.Local().Format("2 January 2006 • 15:04 MST")
	}
	return rv, err
}

const (
	displayWidth = 80
	indentWidth  = 4
	ellipsis     = "..."
)

func capText(txt string) string {
	if len(txt) > (displayWidth - indentWidth) {
		return fmt.Sprintf("%s%s", txt[:displayWidth-indentWidth-len(ellipsis)], ellipsis)
	}
	return txt
}
