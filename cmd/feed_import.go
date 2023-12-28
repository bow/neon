// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/entity"
)

func newFeedImportCommand() *cobra.Command {
	var name = "import"

	command := cobra.Command{
		Use:     fmt.Sprintf("%s [input]", name),
		Args:    cobra.MaximumNArgs(1),
		Aliases: makeAlias(name),
		Short:   "Import feeds from OPML",
		Example: fmt.Sprintf(`  - Import from stdin  : cat feeds.opml | %[1]s feed import
  - Import from a file : %[1]s feed import feeds.opml`, internal.AppName()),

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				err      error
				contents []byte
			)
			if len(args) == 0 {
				contents, err = io.ReadAll(os.Stdin)
			} else {
				switch fn := args[0]; fn {
				case "-", "/dev/stdin":
					contents, err = io.ReadAll(os.Stdin)
				default:
					contents, err = os.ReadFile(fn)
				}
			}
			if err != nil {
				return err
			}

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			sub, err := entity.NewSubscriptionFromRawOPML(contents)
			if err != nil {
				return fmt.Errorf("failed to parse OPML document: %w", err)
			}

			nproc, nimp, err := db.ImportSubscription(cmd.Context(), sub)
			if err != nil {
				return err
			}
			log.Info().
				Int("num_processed", nproc).
				Int("num_imported", nimp).
				Msg("finished feed import")

			return nil
		},
	}

	return &command
}
