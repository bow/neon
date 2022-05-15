package cmd

import (
	"os"

	"github.com/bow/courier/logging"
	"github.com/bow/courier/version"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:               version.AppName(),
	Short:             "RSS reader suite",
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Expose args as proper flags.
		return logging.Init("info", logging.PrettyConsoleStyle, os.Stderr)
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
