// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal"
)

// newVersionCmd creates a new 'version' subcommand.
func newVersionCmd() *cobra.Command {

	versionCmd := cobra.Command{
		Use:   "version",
		Short: "Show the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			bi, ok := debug.ReadBuildInfo()
			if !ok {
				return fmt.Errorf("could not read build info")
			}

			var os, arch = "?", "?"
			for _, setting := range bi.Settings {
				switch setting.Key {
				case "GOOS":
					os = setting.Value
				case "GOARCH":
					arch = setting.Value
				}
			}

			showVersionAttr("App", internal.AppName())
			showVersionAttr("Version", internal.Version())
			showVersionAttr("Git commit", internal.GitCommit())
			showVersionAttr("Build time", internal.BuildTime())
			showVersionAttr("OS/Arch", fmt.Sprintf("%s/%s", os, arch))
			showVersionAttr("Go version", bi.GoVersion)

			return nil
		},
	}

	return &versionCmd
}

// showVersionAttr prints a version-related key-value pair to stdout.
func showVersionAttr(key, value string) {
	fmt.Printf("%-11s: %s\n", key, value)
}
