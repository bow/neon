// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bow/lens/internal/server"
)

// newServerCommand creates a new 'server' subcommand along with its command-line flags.
func newServerCommand() *cobra.Command {

	var (
		name        = "server"
		v           = newViper(name)
		defaultAddr = "127.0.0.1:5151"
	)

	const (
		quietKey = "quiet"
		addrKey  = "addr"
	)

	command := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Start a gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {

			if !v.GetBool(quietKey) {
				showBanner(cmd.OutOrStdout())
			}

			dbPath, err := resolveDBPath(v.GetString(dbPathKey))
			if err != nil {
				return err
			}

			server, err := server.NewBuilder().
				Context(cmd.Context()).
				Address(qualifyAddr(v.GetString(addrKey))).
				StorePath(dbPath).
				Build()

			if err != nil {
				return err
			}

			return server.Serve(cmd.Context())
		},
	}

	flags := command.Flags()

	flags.BoolP(quietKey, "q", false, "hide startup banner")
	flags.StringP(addrKey, "a", defaultAddr, "listening address")
	flags.StringP(dbPathKey, "d", defaultDBPath, "data store location")

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	command.AddCommand(newServerShowProtoCommand())

	return &command
}

// qualifyAddr ensures the specified address has either a 'tcp' or 'file' protocol. If the
// input has no protocol prefix, 'tcp' is assumed.
func qualifyAddr(addr string) string {
	if !server.IsTCPAddr(addr) && !server.IsFileAddr(addr) {
		addr = fmt.Sprintf("tcp://%s", addr)
	}
	return strings.ToLower(addr)
}
