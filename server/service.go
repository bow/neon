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

// GetVersion satisfies the service API.
func (svc *service) GetVersion(
	_ context.Context,
	_ *proto.GetVersionRequest,
) (*proto.GetVersionResponse, error) {

	rsp := proto.GetVersionResponse{
		Name:      version.AppName(),
		Version:   version.Version(),
		GitCommit: version.GitCommit(),
		BuildTime: version.BuildTime(),
	}

	return &rsp, nil
}
