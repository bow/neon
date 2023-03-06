// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"errors"
	"fmt"

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

			ch := str.PullFeeds(cmd.Context())

			var (
				errs []error
				n    int
			)
			for pr := range ch {
				if err := pr.Error(); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w", pr.URL(), err))
					log.Error().
						Str("url", pr.URL()).
						Str("title", pr.Feed().Title).
						Msg("Feed pull failed")
				} else {
					n++
					log.Info().
						Str("url", pr.URL()).
						Str("title", pr.Feed().Title).
						Msg("Feed pull OK")
				}
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			log.Info().Int("num_updated", n).Msgf("Finished pulling feeds")

			return nil
		},
	}

	return &exportCmd
}
