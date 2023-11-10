// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"net"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/bow/iris/api"
	"github.com/bow/iris/internal"
)

func defaultTestServerBuilder(t *testing.T) *Builder {
	t.Helper()

	return NewBuilder().
		Address("file://" + t.TempDir() + "/iris.socket").
		Store(internal.NewMockFeedStore(gomock.NewController(t))).
		Logger(zerolog.Nop())
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

func (tcb *testClientBuilder) ServerStore(str internal.FeedStore) *testClientBuilder {
	tcb.serverBuilder = tcb.serverBuilder.Store(str)
	return tcb
}

func (tcb *testClientBuilder) Build() api.IrisClient {
	tcb.t.Helper()

	t := tcb.t

	// TODO: Avoid global states like this.
	zerolog.SetGlobalLevel(zerolog.Disabled)

	b := tcb.serverBuilder
	if b == nil {
		b = defaultTestServerBuilder(t)
	}
	srv, addr := newTestServer(t, b)

	dialOpts := tcb.dialOpts
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, conn := newTestClient(t, addr, dialOpts...)

	t.Cleanup(
		func() {
			require.NoError(t, conn.Close())
			srv.Stop()
		},
	)

	return client
}

func newTestServer(t *testing.T, b *Builder) (*server, net.Addr) {
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

	return srv, srv.lis.Addr()
}

func newTestClient(
	t *testing.T,
	addr net.Addr,
	opts ...grpc.DialOption,
) (api.IrisClient, *grpc.ClientConn) {
	t.Helper()

	dialer := func(_ context.Context, rawAddr string) (net.Conn, error) {
		return net.Dial(addr.Network(), rawAddr)
	}
	opts = append(opts, grpc.WithContextDialer(dialer))
	conn, err := grpc.Dial(addr.String(), opts...)
	require.NoError(t, err)
	client := api.NewIrisClient(conn)

	return client, conn
}

// setupServerTest is a shortcut method for creating server tests through a client.
func setupServerTest(t *testing.T) (api.IrisClient, *internal.MockFeedStore) {
	t.Helper()

	str := internal.NewMockFeedStore(gomock.NewController(t))
	clb := newTestClientBuilder(t).ServerStore(str)

	return clb.Build(), str
}

func TestServerBuilderErrInvalidAddr(t *testing.T) {
	b := NewBuilder().Address("invalid")
	srv, err := b.Build()
	assert.Nil(t, srv)
	assert.EqualError(t, err, "unexpected address type: invalid")
}
