package internal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bow/courier/api"
)

// service implements the Courier service API.
type service struct {
	api.UnimplementedCourierServer

	store  FeedStore
	parser FeedParser
}

func setupService(grpcs *grpc.Server, store FeedStore, parser FeedParser) *service {
	svc := service{store: store, parser: parser}
	api.RegisterCourierServer(grpcs, &svc)
	return &svc
}

// EditFeed satisfies the service API.
func (svc *service) EditFeed(
	_ context.Context,
	_ *api.EditFeedRequest,
) (*api.EditFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// ListFeeds satisfies the service API.
func (svc *service) ListFeeds(
	_ context.Context,
	_ *api.ListFeedsRequest,
) (*api.ListFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// DeleteFeeds satisfies the service API.
func (svc *service) DeleteFeeds(
	_ context.Context,
	_ *api.DeleteFeedsRequest,
) (*api.DeleteFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// PollFeeds satisfies the service API.
func (svc *service) PollFeeds(_ api.Courier_PollFeedsServer) error {
	return status.Errorf(codes.Unimplemented, "unimplemented")
}

// EditEntry satisfies the service API.
func (svc *service) EditEntry(
	_ context.Context,
	_ *api.EditEntryRequest,
) (*api.EditEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// ExportOPML satisfies the service API.
func (svc *service) ExportOPML(
	_ context.Context,
	_ *api.ExportOPMLRequest,
) (*api.ExportOPMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// ImportOPML satisfies the service API.
func (svc *service) ImportOPML(
	_ context.Context,
	_ *api.ImportOPMLRequest,
) (*api.ImportOPMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
