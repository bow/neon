// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
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

			out := cmd.OutOrStdout()
			showVersionAttr(out, "App", internal.AppName())
			showVersionAttr(out, "Version", internal.Version())
			showVersionAttr(out, "Git commit", internal.GitCommit())
			showVersionAttr(out, "Build time", internal.BuildTime())
			showVersionAttr(out, "OS/Arch", fmt.Sprintf("%s/%s", os, arch))
			showVersionAttr(out, "Go version", bi.GoVersion)

			return nil
		},
	}

	return &versionCmd
}

// showVersionAttr prints a version-related key-value pair to the given writer.
func showVersionAttr(w io.Writer, key, value string) {
	fmt.Fprintf(w, "%-11s: %s\n", key, value)
}
