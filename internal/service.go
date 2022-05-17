package internal

import (
	"context"

	"google.golang.org/grpc"

	"github.com/bow/courier/api"
)

// service implements the Courier service API.
type service struct {
	api.UnimplementedCourierServer
}

func setupService(grpcs *grpc.Server) *service {
	svc := service{}
	api.RegisterCourierServer(grpcs, &svc)
	return &svc
}

// GetInfo satisfies the service API.
func (svc *service) GetInfo(
	_ context.Context,
	_ *api.GetInfoRequest,
) (*api.GetInfoResponse, error) {

	rsp := api.GetInfoResponse{
		Name:      AppName(),
		Version:   Version(),
		GitCommit: GitCommit(),
		BuildTime: BuildTime(),
	}

	return &rsp, nil
}
