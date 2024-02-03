// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bow/neon/internal/entity"
)

func newFeedListCommand() *cobra.Command {
	const name = "list"

	command := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "List feeds",

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			feeds, err := db.ListFeeds(cmd.Context(), nil)
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

func fmtFeed(feed *entity.Feed) string {
	var (
		sb  strings.Builder
		cat = func(format string, a ...any) { fmt.Fprintf(&sb, format, a...) }
	)

	fmtTags := func(tags []string) string {
		if len(tags) > 0 {
			return fmt.Sprintf("#%s", strings.Join(tags, " #"))
		}
		return "-"
	}

	kv := []*struct {
		k, v string
	}{
		{"FeedID", fmt.Sprintf("%d", feed.ID)},
		{"Last pulled", fmtTime(feed.LastPulled)},
		{"Updated", fmtOrEmpty(feed.Updated)},
		{"Unread", fmt.Sprintf("%d/%d", feed.NumEntriesUnread(), feed.NumEntriesTotal())},
		{"URL", derefOrEmpty(feed.SiteURL)},
		{"Tags", fmtTags(feed.Tags)},
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
