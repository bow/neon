package cmd

import (
	"fmt"
	"os"

	"github.com/bow/courier/logging"
	"github.com/bow/courier/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	logLevelKey = "log-level"
	logStyleKey = "log-style"
)

var rootCmd = cobra.Command{
	Use:               version.AppName(),
	Short:             "RSS reader suite",
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel := viper.GetString(logLevelKey)

		var ls logging.Style
		switch rls := viper.GetString(logStyleKey); rls {
		case "pretty":
			ls = logging.PrettyConsoleStyle
		case "json":
			ls = logging.JSONStyle
		default:
			return fmt.Errorf("invalid %s value: '%s'", logStyleKey, rls)
		}

		return logging.Init(logLevel, ls, os.Stderr)
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	pflags := rootCmd.PersistentFlags()

	pflags.StringP(logLevelKey, "l", "info", "logging level")
	viper.BindPFlag(logLevelKey, pflags.Lookup(logLevelKey))

	pflags.String(logStyleKey, "pretty", "logging style")
	viper.BindPFlag(logStyleKey, pflags.Lookup(logStyleKey))
}
