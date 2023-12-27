// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/opml"
)

func newFeedExportCommand() *cobra.Command {
	var name = "export"

	command := cobra.Command{
		Use:     fmt.Sprintf("%s [output]", name),
		Args:    cobra.MaximumNArgs(1),
		Aliases: makeAlias(name),
		Short:   "Export feeds to OPML",
		Example: fmt.Sprintf(`  - Export to stdout : %[1]s feed export
  - Export to a file : %[1]s feed export feeds.opml`, internal.AppName()),

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				err  error
				dest io.Writer
			)
			if len(args) == 0 {
				dest = os.Stdout
			} else {
				switch fn := args[0]; fn {
				case "-", "", "/dev/stdout":
					dest = os.Stdout
				default:
					var fh *os.File
					fh, err = os.Create(fn)
					if err != nil {
						return err
					}
					dest = fh
				}
			}

			db, err := dbFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			sub, err := db.ExportSubscription(cmd.Context(), nil)
			if err != nil {
				return err
			}

			payload, err := sub.Export()
			if err != nil {
				if errors.Is(err, opml.ErrEmptyDocument) {
					return fmt.Errorf("nothing to export")
				}
			}

			_, err = dest.Write(payload)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return &command
}
