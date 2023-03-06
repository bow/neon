// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bow/iris/api"
)

// newServerProtoCmd creates a new 'server' subcommand along with its command-line flags.
func newServerProtoCmd() *cobra.Command {

	var name = "proto"

	protoCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Show the server proto file",

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", api.Proto())
			return nil
		},
	}

	return &protoCmd
}
