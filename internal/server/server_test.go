// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/datastore"
)

func defaultTestServerBuilder(t *testing.T) *Builder {
	t.Helper()

	return NewBuilder().
		Address("tcp://:0").
		Datastore(NewMockDatastore(gomock.NewController(t)))
}

type testClientBuilder struct {
	t             *testing.T
	serverBuilder *Builder
	dialOpts      []grpc.DialOption
}

func newTestClientBuilder(t *testing.T) *testClientBuilder {
	t.Helper()
	return &testClientBuilder{t: t, serverBuilder: defaultTestServerBuilder(t)}
}

func (tcb *testClientBuilder) DialOpts(opts ...grpc.DialOption) *testClientBuilder {
	tcb.dialOpts = opts
	return tcb
}

func (tcb *testClientBuilder) ServerDatastore(ds datastore.Datastore) *testClientBuilder {
	tcb.serverBuilder = tcb.serverBuilder.Datastore(ds)
	return tcb
}

func (tcb *testClientBuilder) Build() api.NeonClient {
	tcb.t.Helper()

	t := tcb.t

	SetLogger(zerolog.Nop())

	b := tcb.serverBuilder
	if b == nil {
		b = defaultTestServerBuilder(t)
	}
	srv := newTestServer(t, b)

	dialOpts := tcb.dialOpts
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, conn := newTestClient(t, srv.Addr(), dialOpts...)

	t.Cleanup(
		func() {
			require.NoError(t, conn.Close())
			srv.Stop()
		},
	)

	return client
}

func newTestServer(t *testing.T, b *Builder) *Server {
	t.Helper()

	r := require.New(t)

	srv, err := b.Build()
	r.NoError(err)

	go func() {
		r.NoError(srv.Serve(context.Background()))
	}()

	var (
		req     = grpc_health_v1.HealthCheckRequest{Service: srv.ServiceName()}
		freq    = 100 * time.Millisecond
		ticker  = time.NewTicker(freq)
		timeout = 5 * time.Second
		timer   = time.NewTimer(timeout)
	)
	defer timer.Stop()
	defer ticker.Stop()

startwait:
	for {
		select {
		case <-ticker.C:
			rsp, err := srv.healthSvc.Check(context.Background(), &req)
			if err != nil {
				t.Fatalf("service health check: %s", err)
			}
			if rsp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
				break startwait
			}
		case <-timer.C:
			t.Fatalf("server startup exceeded maximum time of %s", timeout)
		}
	}

	return srv
}

func newTestClient(
	t *testing.T,
	addr net.Addr,
	opts ...grpc.DialOption,
) (api.NeonClient, *grpc.ClientConn) {
	t.Helper()

	dialer := func(_ context.Context, rawAddr string) (net.Conn, error) {
		return net.Dial(addr.Network(), rawAddr)
	}
	opts = append(opts, grpc.WithContextDialer(dialer))
	conn, err := grpc.NewClient(addr.String(), opts...)
	require.NoError(t, err)
	client := api.NewNeonClient(conn)

	return client, conn
}

// setupServerTest is a shortcut method for creating server tests through a client.
func setupServerTest(t *testing.T) (api.NeonClient, *MockDatastore) {
	t.Helper()

	ds := NewMockDatastore(gomock.NewController(t))
	clb := newTestClientBuilder(t).ServerDatastore(ds)

	return clb.Build(), ds
}

func TestServerBuilderErrInvalidAddr(t *testing.T) {
	b := NewBuilder().Address("invalid")
	srv, err := b.Build()
	assert.Nil(t, srv)
	assert.EqualError(t, err, "unexpected address type: invalid")
}
