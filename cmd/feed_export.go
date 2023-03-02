package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/store/opml"
)

func newFeedExportCmd() *cobra.Command {
	var (
		name        = "export"
		exportViper = newViper(name)
	)

	exportCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Export feeds to OPML",
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := storeFromCtx(cmd)
			if err != nil {
				return err
			}

			contents, err := str.ExportOPML(cmd.Context(), nil)
			if err != nil {
				if errors.Is(err, opml.ErrEmptyDocument) {
					return fmt.Errorf("nothing to export")
				}
				return err
			}

			var (
				out  = exportViper.GetString(exportOutKey)
				dest io.Writer
			)
			switch out {
			case "-", "", "/dev/stdout":
				dest = os.Stdout
			default:
				var fh *os.File
				fh, err = os.Create(out)
				if err != nil {
					return err
				}
				dest = fh
			}

			_, err = dest.Write(contents)
			if err != nil {
				return err
			}

			return nil
		},
	}

	flags := exportCmd.Flags()

	flags.StringP(exportOutKey, "o", "/dev/stdout", "exported path file")
	_ = exportViper.BindPFlag(exportOutKey, flags.Lookup(exportOutKey))

	return &exportCmd
}
