// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal"
)

const (
	logLevelKey = "log-level"
	logStyleKey = "log-style"
	quietKey    = "quiet"
	inDockerKey = "in-docker"
)

// Execute runs the root command.
func Execute() error {
	return newCmd().Execute()
}

// newCmd creates a new root command and sets up its command-line flags.
func newCmd() *cobra.Command {

	var cmdViper = newViper("")

	root := cobra.Command{
		Use:               internal.AppName(),
		Short:             "Feed reader suite",
		SilenceUsage:      true,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			logLevel := cmdViper.GetString(logLevelKey)

			var ls internal.LogStyle
			switch rls := cmdViper.GetString(logStyleKey); rls {
			case "pretty":
				ls = internal.PrettyLogStyle
			case "json":
				ls = internal.JSONLogStyle
			default:
				return fmt.Errorf("invalid %s value: '%s'", logStyleKey, rls)
			}

			err := internal.InitGlobalLog(logLevel, ls, os.Stderr)
			if err != nil {
				return err
			}
			if !cmdViper.GetBool(inDockerKey) {
				internal.SetLogPID()
			}

			if !cmdViper.GetBool(quietKey) {
				showBanner()
			}

			return nil
		},
	}

	pflags := root.PersistentFlags()

	pflags.BoolP(quietKey, "q", false, "hide startup banner")
	_ = cmdViper.BindPFlag(quietKey, pflags.Lookup(quietKey))

	pflags.StringP(logLevelKey, "l", "info", "logging level")
	_ = cmdViper.BindPFlag(logLevelKey, pflags.Lookup(logLevelKey))

	pflags.String(logStyleKey, "pretty", "logging style")
	_ = cmdViper.BindPFlag(logStyleKey, pflags.Lookup(logStyleKey))

	pflags.BoolP(inDockerKey, "", false, "indicate if execution is inside docker")
	_ = cmdViper.BindPFlag(inDockerKey, pflags.Lookup(inDockerKey))

	root.AddCommand(newVersionCmd())
	root.AddCommand(newServeCmd())

	return &root
}

func showBanner() {
	fmt.Printf(`    ____       _
   /  _/_____ (_)_____
   / / / ___// // ___/
 _/ / / /   / /(__  )
/___//_/   /_//____/

`)
}
