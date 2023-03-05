package cmd

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newFeedImportCmd() *cobra.Command {
	var (
		name        = "import"
		importViper = newViper(name)
	)

	importCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Import feeds from OPML",
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := dataStoreFromCmdCtx(cmd)
			if err != nil {
				return err
			}

			var (
				in       = importViper.GetString(importInKey)
				contents []byte
			)
			switch in {
			case "-", "/dev/stdin":
				contents, err = io.ReadAll(os.Stdin)
			default:
				contents, err = os.ReadFile(in)
			}
			if err != nil {
				return err
			}

			n, err := str.ImportOPML(cmd.Context(), contents)
			if err != nil {
				return err
			}
			// TODO: Only show newly imported feeds.
			log.Info().Msgf("imported %d feed(s)", n)

			return nil
		},
	}

	flags := importCmd.Flags()

	// TODO: Set as positional argument.
	flags.StringP(importInKey, "i", "", "path to file to import (required)")
	_ = importCmd.MarkFlagRequired(importInKey)
	_ = importViper.BindPFlag(importInKey, flags.Lookup(importInKey))

	return &importCmd
}
