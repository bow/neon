// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/bow/neon/internal/entity"
)

func newFeedAddCommand() *cobra.Command {

	const (
		name       = "add"
		titleKey   = "title"
		descKey    = "desc"
		starKey    = "star"
		tagKey     = "tag"
		timeoutKey = "timeout"
	)
	var v = newViper(name)

	command := cobra.Command{
		Use:     fmt.Sprintf("%s [input]", name),
		Args:    cobra.ExactArgs(1),
		Aliases: makeAlias(name),
		Short:   "Add a new feed",

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			url := args[0]

			var title *string
			if value := v.GetString(titleKey); value != "" {
				title = &value
			}

			var desc *string
			if value := v.GetString(descKey); value != "" {
				desc = &value
			}

			var isStarred *bool
			if value := v.GetBool(starKey); value {
				isStarred = &value
			}

			var tags []string
			if value := v.GetStringSlice(tagKey); len(value) > 0 {
				tags = value
			}

			var pullTimeout *time.Duration
			if value := v.GetDuration(timeoutKey); value > 0 {
				pullTimeout = &value
			}

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			feed, added, err := db.AddFeed(
				cmd.Context(),
				url,
				title,
				desc,
				tags,
				isStarred,
				pullTimeout,
			)
			if err != nil {
				return err
			}

			logAddResult(feed, added)

			return nil
		},
	}

	flags := command.Flags()

	flags.StringP(titleKey, "t", "", "feed title")
	flags.String(descKey, "", "feed description")
	flags.Bool(starKey, false, "star the feed")
	flags.StringArray(tagKey, nil, "feed tags")
	flags.Duration(timeoutKey, 20*time.Second, "timeout for adding the feed")

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	return &command
}

func logAddResult(feed *entity.Feed, added bool) {

	var msg string
	if added {
		msg = "added feed"
	} else {
		msg = "refreshed feed"
	}

	l := log.Info()
	if feed.FeedURL != "" {
		l = l.Str("feed_url", feed.FeedURL)
	}
	if feed.Title != "" {
		l = l.Str("title", feed.Title)
	}
	if feed.SiteURL != nil {
		l = l.Str("site_url", *feed.SiteURL)
	}
	if feed.IsStarred {
		l = l.Bool("starred", feed.IsStarred)
	}
	if len(feed.Tags) > 0 {
		l = l.Strs("tags", feed.Tags)
	}

	l.Msg(msg)
}
