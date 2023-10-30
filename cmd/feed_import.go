// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newFeedImportCommand() *cobra.Command {
	var name = "import"

	command := cobra.Command{
		Use:     fmt.Sprintf("%s [input]", name),
		Args:    cobra.MaximumNArgs(1),
		Aliases: makeAlias(name),
		Short:   "Import feeds from OPML",
		Example: `  - Import from stdin  : cat feeds.opml | iris feed import
  - Import from a file : iris feed import feeds.opml`,

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

			str, err := dataStoreFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			nproc, nimp, err := str.ImportOPML(cmd.Context(), contents)
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
