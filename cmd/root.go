// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/bow/iris/internal"
)

const quietKey = "quiet"

// New creates a new command along with its command-line flags.
func New() *cobra.Command {

	var cmdViper = newViper("")

	root := cobra.Command{
		Use:                internal.AppName(),
		Short:              "Feed reader suite",
		SilenceUsage:       true,
		SilenceErrors:      true,
		DisableSuggestions: true,
		CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if !cmdViper.GetBool(quietKey) {
				showBanner(cmd.OutOrStdout())
			}

			caser := cases.Title(language.English)

			log.Debug().
				Str("version", internal.Version()).
				Int("pid", os.Getpid()).
				Bool("in_docker", inDocker()).
				Msgf("starting %s", caser.String(internal.AppName()))

			return nil
		},
	}

	pflags := root.PersistentFlags()

	pflags.BoolP(quietKey, "q", false, "hide startup banner")
	_ = cmdViper.BindPFlag(quietKey, pflags.Lookup(quietKey))

	root.AddCommand(newVersionCmd())
	root.AddCommand(newServerCmd())

	return &root
}

// showBanner prints the application banner to the given writer.
func showBanner(w io.Writer) {
	fmt.Fprintf(
		w,
		`    ____       _
   /  _/_____ (_)_____
   / / / ___// // ___/
 _/ / / /   / /(__  )
/___//_/   /_//____/

`)
}

func inDocker() bool {
	_, errStat := os.Stat("/.dockerenv")
	return errStat == nil
}
