// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/entity"
)

func newFeedPullCommand() *cobra.Command {
	var name = "pull"
	const numMaxIDs = 500

	command := cobra.Command{
		Use:     fmt.Sprintf("%s [FEED-ID...]", name),
		Aliases: makeAlias(name),
		Short:   "Pull feed entries",

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			nargs := len(args)
			if nargs > numMaxIDs {
				return fmt.Errorf("number of specified feeds exceeds %d", numMaxIDs)
			}

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			rawIDs := internal.Dedup(args)
			ids, err := entity.ToFeedIDs(rawIDs)
			if err != nil {
				return err
			}

			var (
				errs []error
				n    int
				s    = newPullSpinner(rawIDs)
				ch   = db.PullFeeds(cmd.Context(), ids)
			)

			s.Start()
			defer s.Stop()
			for pr := range ch {
				if err := pr.Error(); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w", pr.URL(), err))
				} else {
					n++
				}
			}
			s.Stop()

			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			log.Info().Int("num_pulled", n).Msgf("Finished pulling feeds")

			return nil
		},
	}

	return &command
}

func newPullSpinner(rawIDs []string) *spinner.Spinner {
	var msg string
	if nids := len(rawIDs); nids == 0 {
		msg = "Pulling all feeds..."
	} else {
		if nids == 1 {
			msg = fmt.Sprintf("Pulling feeds with ID=%v...", rawIDs[0])
		} else {
			var elems []string
			for _, rid := range rawIDs {
				elems = append(elems, fmt.Sprintf("%v", rid))
			}
			msg = fmt.Sprintf("Pulling feeds with IDs=[%s]...", strings.Join(elems, ","))
		}
	}
	return spinner.New(
		spinnerChars,
		75*time.Millisecond,
		spinner.WithColor("cyan"),
		spinner.WithSuffix(" "+bold(msg)),
	)
}

var spinnerChars = []string{
	"█████",
	"▒████",
	"▒▒███",
	"▒▒▒██",
	"█▒▒▒█",
	"██▒▒▒",
	"███▒▒",
	"████▒",
	"█████",
	"████▒",
	"███▒▒",
	"██▒▒▒",
	"█▒▒▒█",
	"▒▒▒██",
	"▒▒███",
	"▒████",
}

func bold(s any) string {
	return fmt.Sprintf("\x1b[1m%v\x1b[0m", s)
}
