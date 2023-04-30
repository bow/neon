// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newFeedPullCmd() *cobra.Command {
	var name = "pull"

	exportCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Pull feed entries",
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := dataStoreFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			var (
				errs []error
				n    int
				ch   = str.PullFeeds(cmd.Context())
				s    = newSpinner()
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
			log.Info().Int("num_updated", n).Msgf("Finished pulling feeds")

			return nil
		},
	}

	return &exportCmd
}

func newSpinner() *spinner.Spinner {
	return spinner.New(
		spinnerChars,
		75*time.Millisecond,
		spinner.WithColor("cyan"),
		spinner.WithSuffix(" "+bold("Pulling feeds...")),
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
