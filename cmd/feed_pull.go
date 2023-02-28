package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func newFeedPullCmd() *cobra.Command {
	var (
		name        = "pull"
		exportViper = newViper(name)
	)

	exportCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Pull feed entries",
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := storeFromCtx(cmd)
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
				} else {
					n++
				}
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated %d feeds\n", n)

			return nil
		},
	}

	flags := exportCmd.Flags()

	flags.StringP(exportOutKey, "o", "/dev/stdout", "exported path file")
	_ = exportViper.BindPFlag(exportOutKey, flags.Lookup(exportOutKey))

	return &exportCmd
}
