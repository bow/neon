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

// newServerCmd creates a new 'server' subcommand along with its command-line flags.
func newServerCmd() *cobra.Command {

	var (
		name        = "server"
		serverViper = newViper(name)
		defaultAddr = "$XDG_RUNTIME_DIR/iris/iris.socket"
	)

	const (
		quietKey = "quiet"
		addrKey  = "addr"
	)

	serverCmd := cobra.Command{
		Use:     name,
		Aliases: makeAlias(name),
		Short:   "Start a gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {

			if !serverViper.GetBool(quietKey) {
				showBanner(cmd.OutOrStdout())
			}

			dbPath, err := resolveDBPath(serverViper.GetString(dbPathKey))
			if err != nil {
				return err
			}

			addr, err := resolveUDSAddr(serverViper.GetString(addrKey))
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

	flags := serverCmd.Flags()

	flags.BoolP(quietKey, "q", false, "hide startup banner")
	_ = serverViper.BindPFlag(quietKey, flags.Lookup(quietKey))

	flags.StringP(addrKey, "a", defaultAddr, "listening address")
	_ = serverViper.BindPFlag(addrKey, flags.Lookup(addrKey))

	flags.StringP(dbPathKey, "d", defaultDBPath, "data store location")
	_ = serverViper.BindPFlag(dbPathKey, flags.Lookup(dbPathKey))

	return &serverCmd
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
