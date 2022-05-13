package server

import (
	"errors"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	lis net.Listener

	grpcServer *grpc.Server
}

func (s *server) Serve() error {
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
	addr string
}

func NewBuilder() *Builder {
	builder := Builder{}
	return &builder
}

func (b *Builder) Address(addr string) *Builder {
	b.addr = addr
	return b
}

func (b *Builder) Build() (*server, error) {
	lis, err := net.Listen("tcp", b.addr)
	if err != nil {
		return nil, err
	}

	grpcs := grpc.NewServer()
	setupService(grpcs)

	srv := server{grpcServer: grpcs, lis: lis}

	return &srv, nil
}
