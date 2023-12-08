// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/server"
)

// newServerCommand creates a new 'server' subcommand along with its command-line flags.
func newServerCommand() *cobra.Command {

	var (
		name        = "server"
		v           = newViper(name)
		defaultAddr = "$XDG_RUNTIME_DIR/iris/server.socket"
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

			addr, err := resolveUDSAddr(v.GetString(addrKey))
			if err != nil {
				return err
			}

			server, err := server.NewBuilder().
				Address(addr).
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

// resolveUDSAddr attempts to resolve the filesystem path to a Unix domain socket exposing
// the running application.
func resolveUDSAddr(addr string) (string, error) {
	var err error

	// return string unchanged if tcp is requested.
	if strings.HasPrefix(strings.ToLower(addr), "tcp") {
		return addr, nil
	}

	xdgDir := "$XDG_RUNTIME_DIR/"
	if strings.HasPrefix(addr, xdgDir) {
		rel := strings.TrimPrefix(addr, xdgDir)
		addr, err = xdg.RuntimeFile(rel)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("file://%s", addr), nil
}
