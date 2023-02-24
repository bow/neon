package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	exportOutKey = "export-out"
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
				return err
			}

			var (
				out  = exportViper.GetString(exportOutKey)
				dest io.Writer
			)
			switch out {
			case "-", "":
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

	flags.StringP(exportOutKey, "o", "-", "exported path file")
	_ = exportViper.BindPFlag(exportOutKey, flags.Lookup(exportOutKey))

	return &exportCmd
}
