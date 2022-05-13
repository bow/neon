package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/bow/courier/version"
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

		show("App", version.AppName())
		show("Version", version.Version())
		show("Git commit", version.GitCommit())
		show("Build time", version.BuildTime())
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
