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

			showVersion(
				cmd.OutOrStdout(),
				map[string]string{
					"App":        internal.AppName(),
					"Version":    internal.Version(),
					"Git commit": internal.GitCommit(),
					"Build time": internal.BuildTime(),
					"OS/Arch":    fmt.Sprintf("%s/%s", os, arch),
					"Go version": bi.GoVersion,
				},
			)

			return nil
		},
	}

	return &versionCmd
}

// showVersion prints version-related information to the given writer.
func showVersion(w io.Writer, info map[string]string) {
	for k, v := range info {
		fmt.Fprintf(w, "%-11s: %s\n", k, v)
	}
}
