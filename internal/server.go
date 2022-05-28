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
	quiet      bool
	stopf      func()

	healthSvc *health.Server

	feeds FeedsStore
}

func newServer(
	lis net.Listener,
	grpcServer *grpc.Server,
	feeds FeedsStore,
	quiet bool,
) *server {

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
		quiet:      quiet,
		healthSvc:  health.NewServer(),
		feeds:      feeds,
	}

	return &s
}

func (s *server) ServiceName() string {
	return api.Courier_ServiceDesc.ServiceName
}

func (s *server) Serve() error {
	if !s.quiet {
		fmt.Printf(`   ______                 _
  / ____/___  __  _______(_)__  _____
 / /   / __ \/ / / / ___/ / _ \/ ___/
/ /___/ /_/ / /_/ / /  / /  __/ /
\____/\____/\__,_/_/  /_/\___/_/

`)
	}
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
	addr   string
	feeds  FeedsStore
	logger zerolog.Logger
	quiet  bool
}

func NewServerBuilder() *ServerBuilder {
	builder := ServerBuilder{logger: zerolog.Nop(), quiet: false}
	return &builder
}

func (b *ServerBuilder) Address(addr string) *ServerBuilder {
	b.addr = addr
	return b
}

func (b *ServerBuilder) Store(feeds FeedsStore) *ServerBuilder {
	b.feeds = feeds
	return b
}

func (b *ServerBuilder) Logger(logger zerolog.Logger) *ServerBuilder {
	b.logger = logger
	return b
}

func (b *ServerBuilder) Quiet(quiet bool) *ServerBuilder {
	b.quiet = quiet
	return b
}

func (b *ServerBuilder) Build() (*server, error) {
	lis, err := net.Listen("tcp", b.addr)
	if err != nil {
		return nil, err
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
	setupService(grpcs)

	s := newServer(lis, grpcs, b.feeds, b.quiet)
	healthapi.RegisterHealthServer(grpcs, s.healthSvc)
	reflection.Register(grpcs)

	return s, nil
}
