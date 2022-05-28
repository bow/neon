package cmd

import (
	"fmt"
	"os"

	"github.com/bow/courier/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	logLevelKey = "log-level"
	logStyleKey = "log-style"
	quietKey    = "quiet"
)

var rootCmd = cobra.Command{
	Use:               internal.AppName(),
	Short:             "Feed reader suite",
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel := viper.GetString(logLevelKey)

		var ls internal.LogStyle
		switch rls := viper.GetString(logStyleKey); rls {
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

		if !viper.GetBool(quietKey) {
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

	pflags.BoolP(quietKey, "q", false, "show banner")
	_ = viper.BindPFlag(quietKey, pflags.Lookup(quietKey))

	pflags.StringP(logLevelKey, "l", "info", "logging level")
	_ = viper.BindPFlag(logLevelKey, pflags.Lookup(logLevelKey))

	pflags.String(logStyleKey, "pretty", "logging style")
	_ = viper.BindPFlag(logStyleKey, pflags.Lookup(logStyleKey))
}

func showBanner() {
	fmt.Printf(`   ______                 _
  / ____/___  __  _______(_)__  _____
 / /   / __ \/ / / / ___/ / _ \/ ___/
/ /___/ /_/ / /_/ / /  / /  __/ /
\____/\____/\__,_/_/  /_/\___/_/

`)
}
