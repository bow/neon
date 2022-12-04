// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"os"

	"github.com/bow/iris/internal"
	"github.com/spf13/cobra"
)

const (
	logLevelKey = "log-level"
	logStyleKey = "log-style"
	quietKey    = "quiet"
	inDockerKey = "in-docker"
)

var rootViper = newViper("")

var rootCmd = cobra.Command{
	Use:               internal.AppName(),
	Short:             "Feed reader suite",
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		logLevel := rootViper.GetString(logLevelKey)

		var ls internal.LogStyle
		switch rls := rootViper.GetString(logStyleKey); rls {
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
		if !rootViper.GetBool(inDockerKey) {
			internal.SetLogPID()
		}

		if !rootViper.GetBool(quietKey) {
			showBanner()
		}

		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	pflags := rootCmd.PersistentFlags()

	pflags.BoolP(quietKey, "q", false, "hide startup banner")
	_ = rootViper.BindPFlag(quietKey, pflags.Lookup(quietKey))

	pflags.StringP(logLevelKey, "l", "info", "logging level")
	_ = rootViper.BindPFlag(logLevelKey, pflags.Lookup(logLevelKey))

	pflags.String(logStyleKey, "pretty", "logging style")
	_ = rootViper.BindPFlag(logStyleKey, pflags.Lookup(logStyleKey))

	pflags.BoolP(inDockerKey, "", false, "indicate if execution is inside docker")
	_ = rootViper.BindPFlag(inDockerKey, pflags.Lookup(inDockerKey))
}

func showBanner() {
	fmt.Printf(`    ____       _
   /  _/_____ (_)_____
   / / / ___// // ___/
 _/ / / /   / /(__  )
/___//_/   /_//____/

`)
}
