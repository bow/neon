package server

import (
	"errors"
	"fmt"
	"net"

	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type server struct {
	lis net.Listener

	grpcServer *grpc.Server
}

func (s *server) Serve() error {
	fmt.Printf(`   ______                 _
  / ____/___  __  _______(_)__  _____
 / /   / __ \/ / / / ___/ / _ \/ ___/
/ /___/ /_/ / /_/ / /  / /  __/ /
\____/\____/\__,_/_/  /_/\___/_/

`)
	log.Info().Msg("starting server")
	err := <-s.start()
	if errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}
	return err
}

func (s *server) start() <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.grpcServer.Serve(s.lis)
	}()
	return ch
}

type Builder struct {
	addr   string
	logger zerolog.Logger
}

func NewBuilder() *Builder {
	builder := Builder{logger: zerolog.Nop()}
	return &builder
}

func (b *Builder) Address(addr string) *Builder {
	b.addr = addr
	return b
}

func (b *Builder) Logger(logger zerolog.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) Build() (*server, error) {
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

	srv := server{grpcServer: grpcs, lis: lis}

	return &srv, nil
}
