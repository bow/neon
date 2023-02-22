package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newFeedListCmd() *cobra.Command {
	var (
		name = "list"
	)

	list := cobra.Command{
		Use:   name,
		Short: "List feeds",
		RunE: func(cmd *cobra.Command, args []string) error {

			dbPath := cmd.Context().Value(ctxKey(dbPathKey))
			fmt.Printf("DB path is: %s\n", dbPath)

			return nil
		},
	}

	return &list
}
