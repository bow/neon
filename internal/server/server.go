package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	"github.com/bow/courier/api"
	"github.com/bow/courier/internal/store"
)

const (
	tcpPrefix  = "tcp://"
	filePrefix = "file://"
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

	healthSvc := health.NewServer()
	healthapi.RegisterHealthServer(grpcServer, healthSvc)

	reflection.Register(grpcServer)

	s := server{
		lis:        lis,
		grpcServer: grpcServer,
		stopf:      func() { funcCh <- struct{}{} },
		healthSvc:  healthSvc,
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

type Builder struct {
	addr      string
	store     store.FeedStore
	storePath string
	parser    store.FeedParser
	logger    zerolog.Logger
}

func NewBuilder() *Builder {
	builder := Builder{logger: log.With().Logger(), parser: gofeed.NewParser()}
	return &builder
}

func (b *Builder) Address(addr string) *Builder {
	b.addr = addr
	return b
}

func (b *Builder) StorePath(path string) *Builder {
	b.storePath = path
	b.store = nil
	return b
}

func (b *Builder) Store(str store.FeedStore) *Builder {
	b.store = str
	b.storePath = ""
	return b
}

func (b *Builder) Parser(parser store.FeedParser) *Builder {
	b.parser = parser
	return b
}

func (b *Builder) Logger(logger zerolog.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) Build() (*server, error) {

	var netw string
	switch addr := b.addr; {
	case isTCPAddr(addr):
		netw = "tcp"
		b.addr = addr[len(tcpPrefix):]
	case isFileAddr(addr):
		netw = "unix"
		b.addr = addr[len(filePrefix):]
	default:
		return nil, fmt.Errorf("unexpected address type: %s", b.addr)
	}

	lis, err := net.Listen(netw, b.addr)
	if err != nil {
		return nil, err
	}

	str := b.store
	if sp := b.storePath; sp != "" {
		if str, err = store.NewSQLite(sp, b.parser); err != nil {
			return nil, fmt.Errorf("server build: %w", err)
		}
	}

	if b.parser == nil {
		b.parser = gofeed.NewParser()
	}

	logger := b.logger.With().
		Str("grpc.version", grpc.Version).
		Logger()

	grpcs := grpc.NewServer(
		middleware.WithUnaryServerChain(
			storeErrorUnaryServerInterceptor,
			tags.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger)),
		),
		middleware.WithStreamServerChain(
			storeErrorStreamServerInterceptor,
			tags.StreamServerInterceptor(),
			logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger)),
		),
	)
	_ = newRPC(grpcs, str)

	s := newServer(lis, grpcs)

	return s, nil
}

func isAddrF(prefix string) func(string) bool {
	return func(addr string) bool {
		return strings.HasPrefix(strings.ToLower(addr), prefix)
	}
}

var (
	isTCPAddr  = isAddrF(tcpPrefix)
	isFileAddr = isAddrF(filePrefix)
)
