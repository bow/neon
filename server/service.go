package server

import (
	"context"

	"google.golang.org/grpc"

	"github.com/bow/courier/proto"
	"github.com/bow/courier/version"
)

// service implements the Courier service API.
type service struct {
	proto.UnimplementedCourierServer
}

func setupService(grpcs *grpc.Server) *service {
	svc := service{}
	proto.RegisterCourierServer(grpcs, &svc)
	return &svc
}

// GetInfo satisfies the service API.
func (svc *service) GetInfo(
	_ context.Context,
	_ *proto.GetInfoRequest,
) (*proto.GetInfoResponse, error) {

	rsp := proto.GetInfoResponse{
		Name:      version.AppName(),
		Version:   version.Version(),
		GitCommit: version.GitCommit(),
		BuildTime: version.BuildTime(),
	}

	return &rsp, nil
}
