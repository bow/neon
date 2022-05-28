package internal

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/bow/courier/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func setupTestServer(t *testing.T) api.CourierClient {
	t.Helper()
	// TODO: Avoid global states like this.
	zerolog.SetGlobalLevel(zerolog.Disabled)

	var (
		b            = NewServerBuilder().Address(":0").Logger(zerolog.Nop())
		srv, addr    = newRunningServer(t, b)
		client, conn = newClient(t, addr)
	)
	t.Cleanup(
		func() {
			require.NoError(t, conn.Close())
			srv.Stop()
		},
	)

	return client
}

func newRunningServer(t *testing.T, b *ServerBuilder) (*server, net.Addr) {
	t.Helper()

	r := require.New(t)

	srv, err := b.Build()
	r.NoError(err)

	go func() {
		r.NoError(srv.Serve())
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

func newClient(
	t *testing.T,
	addr net.Addr,
	opts ...grpc.DialOption,
) (api.CourierClient, *grpc.ClientConn) {
	t.Helper()

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(addr.String(), opts...)
	require.NoError(t, err)
	client := api.NewCourierClient(conn)

	return client, conn
}

func TestServerBuilderErrInvalidAddr(t *testing.T) {
	b := NewServerBuilder().Address("invalid")
	srv, err := b.Build()
	assert.Nil(t, srv)
	assert.EqualError(t, err, "listen tcp: address invalid: missing port in address")
}
