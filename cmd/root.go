package cmd

import "github.com/spf13/cobra"

var rootCmd = cobra.Command{
	Use:               "courier",
	Short:             "RSS reader suite",
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
