// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
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

	rootCmd := cobra.Command{
		Use:                internal.AppName(),
		Short:              "Feed reader suite",
		SilenceUsage:       true,
		SilenceErrors:      true,
		DisableSuggestions: true,
		CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			caser := cases.Title(language.English)

			log.Debug().
				Str("version", internal.Version()).
				Int("pid", os.Getpid()).
				Bool("in_docker", inDocker()).
				Msgf("starting %s", caser.String(internal.AppName()))

			return nil
		},
	}

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newFeedCmd())
	rootCmd.AddCommand(newServerCmd())

	return &rootCmd
}

func inDocker() bool {
	_, errStat := os.Stat("/.dockerenv")
	return errStat == nil
}
