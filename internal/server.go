package internal

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bow/courier/api"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthapi "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct {
	lis        net.Listener
	grpcServer *grpc.Server
	stopf      func()

	healthSvc *health.Server
}

func newServer(lis net.Listener, grpcServer *grpc.Server) *server {

	var (
		funcCh = make(chan struct{}, 1)
		sigCh  = make(chan os.Signal, 1)
	)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer close(funcCh)
		defer close(sigCh)

		select {
		case sig := <-sigCh:
			log.Info().Msgf("stopping server (%s)", sig)
		case <-funcCh:
			log.Info().Msgf("stopping server (function called)")
		}

		grpcServer.GracefulStop()
		log.Info().Msg("server stopped")
	}()

	s := server{
		lis:        lis,
		grpcServer: grpcServer,
		stopf:      func() { funcCh <- struct{}{} },
		healthSvc:  health.NewServer(),
	}

	return &s
}

func (s *server) ServiceName() string {
	return api.Courier_ServiceDesc.ServiceName
}

func (s *server) Serve() error {
	log.Info().Msg("starting server")

	s.healthSvc.SetServingStatus(s.ServiceName(), healthapi.HealthCheckResponse_NOT_SERVING)

	err := <-s.start()
	if errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}
	return err
}

func (s *server) Stop() {
	s.stopf()
}

func (s *server) start() <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		s.healthSvc.Resume()
		ch <- s.grpcServer.Serve(s.lis)
	}()
	log.Info().Msgf("server listening at %s", s.lis.Addr().String())

	return ch
}

type ServerBuilder struct {
	addr      string
	store     FeedsStore
	storePath string
	parser    FeedParser
	logger    zerolog.Logger
}

func NewServerBuilder() *ServerBuilder {
	builder := ServerBuilder{logger: zerolog.Nop()}
	return &builder
}

func (b *ServerBuilder) Address(addr string) *ServerBuilder {
	b.addr = addr
	return b
}

func (b *ServerBuilder) StorePath(path string) *ServerBuilder {
	b.storePath = path
	return b
}

func (b *ServerBuilder) Store(store FeedsStore) *ServerBuilder {
	b.store = store
	return b
}

func (b *ServerBuilder) Parser(parser FeedParser) *ServerBuilder {
	b.parser = parser
	return b
}

func (b *ServerBuilder) Logger(logger zerolog.Logger) *ServerBuilder {
	b.logger = logger
	return b
}

func (b *ServerBuilder) Build() (*server, error) {

	lis, err := net.Listen("tcp", b.addr)
	if err != nil {
		return nil, err
	}

	if b.store != nil && b.storePath != "" {
		return nil, fmt.Errorf("server build: only one of store and storePath may be set")
	}
	if b.store == nil && b.storePath == "" {
		return nil, fmt.Errorf("server build: exactly one of store or storePath must be set")
	}

	store := b.store
	if sp := b.storePath; sp != "" {
		if store, err = newFeedsDB(sp); err != nil {
			return nil, fmt.Errorf("server build: %w", err)
		}
	}

	if b.parser == nil {
		b.parser = gofeed.NewParser()
	}

	grpcs := grpc.NewServer(
		middleware.WithUnaryServerChain(
			tags.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(b.logger)),
		),
		middleware.WithStreamServerChain(
			tags.StreamServerInterceptor(),
			logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(b.logger)),
		),
	)
	setupService(grpcs, store, b.parser)

	s := newServer(lis, grpcs)
	healthapi.RegisterHealthServer(grpcs, s.healthSvc)
	reflection.Register(grpcs)

	return s, nil
}
