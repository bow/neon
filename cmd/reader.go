// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bow/neon/internal/reader"
	"github.com/bow/neon/internal/server"
)

func newReaderCommand() *cobra.Command {
	var (
		name               = "reader"
		v                  = newViper(name)
		defaultStartAddr   = "localhost:0"
		defaultConnectAddr = defaultServerAddr
	)

	const (
		addrKey           = "address"
		connectKey        = "connect"
		connectTimeoutKey = "connect-timeout"
	)

	command := cobra.Command{
		Use:     name,
		Aliases: append(makeAlias(name), []string{"r"}...),
		Short:   "Open the feed reader",
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				err            error
				connectAddr    net.Addr
				connectTimeout time.Duration
				ctx            = cmd.Context()
				dialOpts       = []grpc.DialOption{
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				}
				addr = resolveAddr(v, addrKey, connectKey, defaultConnectAddr, defaultStartAddr)
			)

			if v.GetBool(connectKey) {
				connectAddr, err = makeConnectAddr(addr)
				if err != nil {
					return err
				}
				dialOpts = append(dialOpts, grpc.WithBlock())
				connectTimeout = v.GetDuration(connectTimeoutKey)

			} else {
				server, ierr := makeServer(cmd, v, addr)
				if ierr != nil {
					return ierr
				}

				go func() {
					_ = server.Serve(cmd.Context())
				}()
				defer server.Stop()

				connectAddr = server.Addr()
			}

			rdr, err := reader.NewBuilder(cmd.Context()).
				Context(ctx).
				ConnectTimeout(connectTimeout).
				Address(connectAddr.String()).
				DialOpts(dialOpts...).
				Build()

			if err != nil {
				return err
			}

			return rdr.Start()
		},
	}

	flags := command.Flags()

	flags.StringP(
		addrKey,
		"a",
		"",
		fmt.Sprintf(
			`server address (default "%s" if "-c" is set, localhost with random port otherwise)`,
			defaultConnectAddr,
		),
	)
	flags.BoolP(connectKey, "c", false, "connect to a running server")
	flags.DurationP(
		connectTimeoutKey,
		"t",
		2*time.Second,
		`timeout for initial server connection, ignored if "-c" is unset`,
	)
	flags.StringP(dbPathKey, "d", defaultDBPath, `datastore location, ignored if "-c" is set`)

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	return &command
}

func resolveAddr(
	v *viper.Viper,
	addrKey string,
	connectKey string,
	connectDefault, startDefault string,
) string {
	var (
		addr    string
		connect = v.GetBool(connectKey)
	)

	if v.IsSet(addrKey) {
		addr = v.GetString(addrKey)
	} else {
		if connect {
			addr = connectDefault
		} else {
			addr = startDefault
		}
	}

	return normalizeAddr(addr)
}

func makeConnectAddr(value string) (net.Addr, error) {
	var (
		addr net.Addr
		err  error
	)
	if server.IsTCPAddr(value) {
		addr, err = net.ResolveTCPAddr("tcp", value[len("tcp://"):])
		if err != nil {
			return nil, err
		}
	} else if server.IsFileSystemAddr(value) {
		addr, err = net.ResolveUnixAddr("unix", value[len("file://"):])
		if err != nil {
			return nil, err
		}
	}
	return addr, nil
}
