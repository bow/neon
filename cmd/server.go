// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bow/neon/internal/database"
	"github.com/bow/neon/internal/server"
)

// newServerCommand creates a new 'server' subcommand along with its command-line flags.
func newServerCommand() *cobra.Command {

	var (
		name        = "server"
		v           = newViper(name)
		defaultAddr = defaultServerAddr
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

			database.SetLogger(zlog.Logger)
			server.SetLogger(zlog.Logger)

			if !v.GetBool(quietKey) {
				showBanner(cmd.OutOrStdout())
			}

			srv, err := makeServer(cmd, v, normalizeAddr(v.GetString(addrKey)))
			if err != nil {
				return err
			}

			return srv.Serve(cmd.Context())
		},
	}

	flags := command.Flags()

	flags.BoolP(quietKey, "q", false, "hide startup banner")
	flags.StringP(addrKey, "a", defaultAddr, "listening address")
	flags.StringP(dbPathKey, "d", defaultDBPath, "database location")

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	command.AddCommand(newServerShowProtoCommand())

	return &command
}

func makeServer(cmd *cobra.Command, v *viper.Viper, addr string) (*server.Server, error) {

	dbPath, err := resolveDBPath(v.GetString(dbPathKey))
	if err != nil {
		return nil, err
	}

	srv, err := server.NewBuilder().
		Context(cmd.Context()).
		Address(addr).
		SQLite(dbPath).
		Build()

	return srv, err
}

// normalizeAddr ensures the specified address has either a 'tcp' or 'file' protocol. If the
// input has no protocol prefix, 'tcp' is assumed.
func normalizeAddr(addr string) string {
	if !server.IsTCPAddr(addr) && !server.IsFileSystemAddr(addr) {
		addr = fmt.Sprintf("tcp://%s", addr)
	}
	return strings.ToLower(addr)
}
