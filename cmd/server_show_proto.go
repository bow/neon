// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bow/neon/api"
)

// newServerShowProtoCommand creates a new subcommand for showing the server's protobuf
// interface.
func newServerShowProtoCommand() *cobra.Command {

	const name = "show-proto"

	command := cobra.Command{
		Use:     name,
		Aliases: []string{"sp"},
		Short:   "Show the server protobuf interface",

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", api.Proto())
			return nil
		},
	}

	return &command
}
