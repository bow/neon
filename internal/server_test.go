package internal

import (
	"context"
	"net"
	"path/filepath"
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

func defaultTestServerBuilder(t *testing.T) *ServerBuilder {
	t.Helper()

	return NewServerBuilder().
		Address(":0").
		StorePath(filepath.Join(t.TempDir(), "courier-test.db")).
		Logger(zerolog.Nop())
}

type testClientBuilder struct {
	serverBuilder *ServerBuilder
	dialOpts      []grpc.DialOption
}

func newTestClientBuilder() *testClientBuilder {
	return &testClientBuilder{}
}

func (tcb *testClientBuilder) DialOpts(opts ...grpc.DialOption) *testClientBuilder {
	tcb.dialOpts = opts
	return tcb
}

func (tcb *testClientBuilder) ServerBuilder(b *ServerBuilder) *testClientBuilder {
	tcb.serverBuilder = b
	return tcb
}

func (tcb *testClientBuilder) Build(t *testing.T) api.CourierClient {
	t.Helper()

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

func newTestServer(t *testing.T, b *ServerBuilder) (*server, net.Addr) {
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

func newTestClient(
	t *testing.T,
	addr net.Addr,
	opts ...grpc.DialOption,
) (api.CourierClient, *grpc.ClientConn) {
	t.Helper()

	conn, err := grpc.Dial(addr.String(), opts...)
	require.NoError(t, err)
	client := api.NewCourierClient(conn)

	return client, conn
}

func TestServerBuilderErrInvalidAddr(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "courier-test-err-invalid-addr.db")
	b := NewServerBuilder().Address("invalid").StorePath(storePath)
	srv, err := b.Build()
	assert.Nil(t, srv)
	assert.EqualError(t, err, "listen tcp: address invalid: missing port in address")
}
