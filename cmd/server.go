// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bow/iris/internal/server"
)

const (
	addrKey   = "addr"
	dbNameKey = "db-path"
)

var (
	relDBName     = "iris/iris.db"
	defaultDBName = fmt.Sprintf("$XDG_DATA_HOME/%s", relDBName)

	relUDS      = "iris/iris.socket"
	defaultAddr = fmt.Sprintf("$XDG_RUNTIME_DIR/%s", relUDS)
)

// newServerCmd creates a new 'server' subcommand along with its command-line flags.
func newServerCmd() *cobra.Command {

	var (
		name     = "server"
		cmdViper = newViper(name)
	)

	serverCmd := cobra.Command{
		Use:   name,
		Short: "Start a gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {

			dbPath, err := resolveDBPath(cmdViper)
			if err != nil {
				return err
			}

			addr, err := resolveUDSAddr(cmdViper)
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

	flags.StringP(addrKey, "a", defaultAddr, "listening address")
	_ = cmdViper.BindPFlag(addrKey, flags.Lookup(addrKey))

	flags.StringP(dbNameKey, "d", defaultDBName, "data store location")
	_ = cmdViper.BindPFlag(dbNameKey, flags.Lookup(dbNameKey))

	return &serverCmd
}

// resolveDBPath attempts to resolve the filesystem path to the database.
func resolveDBPath(v *viper.Viper) (dbPath string, err error) {
	dbPath = v.GetString(dbNameKey)
	if dbPath == defaultDBName {
		dbPath, err = xdg.DataFile(relDBName)
		if err != nil {
			return "", err
		}
	}
	return dbPath, nil
}

// resolveUDSAddr attempts to resolve the filesystem path to a Unix domain socket exposing
// the running application.
func resolveUDSAddr(v *viper.Viper) (addr string, err error) {
	addr = v.GetString(addrKey)
	// return string unchanged if tcp is requested.
	if strings.HasPrefix(strings.ToLower(addr), "tcp") {
		return addr, nil
	}
	if addr == defaultAddr {
		addr, err = xdg.RuntimeFile(relUDS)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("file://%s", addr), nil
}
