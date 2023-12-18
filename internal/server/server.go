// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthapi "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/bow/lens/api"
	"github.com/bow/lens/internal"
	"github.com/bow/lens/internal/database"
)

const (
	tcpPrefix  = "tcp://"
	unixPrefix = "unix://"
	filePrefix = "file://"
)

type server struct {
	lis        net.Listener
	grpcServer *grpc.Server
	stopf      func()
	stoppedCh  chan struct{}

	healthSvc *health.Server
}

func newServer(lis net.Listener, grpcServer *grpc.Server, str internal.FeedStore) *server {

	svc := service{store: str}
	api.RegisterLensServer(grpcServer, &svc)

	var (
		funcCh    = make(chan struct{}, 1)
		sigCh     = make(chan os.Signal, 1)
		stoppedCh = make(chan struct{}, 1)
	)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer close(funcCh)
		defer close(sigCh)
		defer close(stoppedCh)

		reason := "unknown"
		select {
		case sig := <-sigCh:
			reason = sig.String()
		case <-funcCh:
			reason = "function called"
		}

		log.Debug().Msg("stopping server")
		grpcServer.GracefulStop()
		log.Info().Msgf("server stopped (%s)", reason)
		stoppedCh <- struct{}{}
	}()

	healthSvc := health.NewServer()
	healthapi.RegisterHealthServer(grpcServer, healthSvc)

	reflection.Register(grpcServer)

	s := server{
		lis:        lis,
		grpcServer: grpcServer,
		stopf:      func() { funcCh <- struct{}{} },
		stoppedCh:  stoppedCh,
		healthSvc:  healthSvc,
	}

	return &s
}

func (s *server) ServiceName() string {
	return api.Lens_ServiceDesc.ServiceName
}

func (s *server) Serve(ctx context.Context) error {
	log.Debug().
		Str("addr", s.lis.Addr().String()).
		Msg("starting server")

	s.healthSvc.SetServingStatus(s.ServiceName(), healthapi.HealthCheckResponse_NOT_SERVING)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-s.start():
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return err
	}
}

func (s *server) Stop() {
	s.stopf()
	<-s.stoppedCh
}

func (s *server) start() <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		s.healthSvc.Resume()
		ch <- s.grpcServer.Serve(s.lis)
	}()
	log.Info().Str("addr", s.lis.Addr().String()).Msgf("server listening")

	return ch
}

type Builder struct {
	ctx       context.Context
	addr      string
	store     internal.FeedStore
	storePath string
	parser    internal.FeedParser
	logger    zerolog.Logger
}

func NewBuilder() *Builder {
	builder := Builder{
		ctx:    context.Background(),
		logger: log.With().Logger(),
		parser: gofeed.NewParser(),
	}
	return &builder
}

func (b *Builder) Context(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
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

func (b *Builder) Store(str internal.FeedStore) *Builder {
	b.store = str
	b.storePath = ""
	return b
}

func (b *Builder) Logger(logger zerolog.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) Build() (*server, error) {

	var netw string
	switch addr := b.addr; {
	case IsTCPAddr(addr):
		netw = "tcp"
		b.addr = addr[len(tcpPrefix):]
	case IsFileAddr(addr):
		netw = "unix"
		b.addr = addr[len(filePrefix):]
	case IsUnixAddr(addr):
		netw = "unix"
		b.addr = addr[len(unixPrefix):]
	default:
		return nil, fmt.Errorf("unexpected address type: %s", b.addr)
	}

	var lc net.ListenConfig
	lis, err := lc.Listen(b.ctx, netw, b.addr)
	if err != nil {
		return nil, err
	}

	str := b.store
	if sp := b.storePath; sp != "" {
		log.Info().Str("path", sp).Msgf("initializing data store")
		if str, err = database.NewSQLiteWithParser(sp, b.parser); err != nil {
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
		grpc.ChainUnaryInterceptor(
			storeErrorUnaryServerInterceptor,
			logging.UnaryServerInterceptor(internal.InterceptorLogger(logger)),
		),
		grpc.ChainStreamInterceptor(
			storeErrorStreamServerInterceptor,
			logging.StreamServerInterceptor(internal.InterceptorLogger(logger)),
		),
	)
	s := newServer(lis, grpcs, str)

	return s, nil
}

func isAddrF(prefix string) func(string) bool {
	return func(addr string) bool {
		return strings.HasPrefix(strings.ToLower(addr), prefix)
	}
}

var (
	IsTCPAddr  = isAddrF(tcpPrefix)
	IsFileAddr = isAddrF(filePrefix)
	IsUnixAddr = isAddrF(unixPrefix)

	IsFileSystemAddr = func(addr string) bool { return IsFileAddr(addr) || IsUnixAddr(addr) }
)
