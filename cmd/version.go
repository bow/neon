// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal"
)

var versionCmd = cobra.Command{
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

		show("App", internal.AppName())
		show("Version", internal.Version())
		show("Git commit", internal.GitCommit())
		show("Build time", internal.BuildTime())
		show("OS/Arch", fmt.Sprintf("%s/%s", os, arch))
		show("Go version", bi.GoVersion)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(&versionCmd)
}

func show(key, value string) {
	fmt.Printf("%-11s: %s\n", key, value)
}
