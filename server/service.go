package server

import (
	"context"

	"github.com/bow/courier/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
