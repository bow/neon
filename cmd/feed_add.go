// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newFeedAddCommand() *cobra.Command {
	var (
		name = "add"
		v    = newViper(name)
	)

	const (
		titleKey = "title"
		descKey  = "desc"
		starKey  = "star"
		tagKey   = "tag"
	)

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

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			feed, err := db.AddFeed(cmd.Context(), url, title, desc, tags, isStarred)
			if err != nil {
				return err
			}
			log.Info().
				Str("feed_url", feed.FeedURL).
				Str("title", feed.Title).
				Msg("added feed")

			return nil
		},
	}

	flags := command.Flags()

	flags.StringP(titleKey, "t", "", "feed title")
	flags.String(descKey, "", "feed description")
	flags.Bool(starKey, false, "star the feed")
	flags.StringArray(tagKey, nil, "feed tags")

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	return &command
}
