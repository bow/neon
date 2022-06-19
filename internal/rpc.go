package internal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bow/courier/api"
)

// rpc implements the Courier rpc API.
type rpc struct {
	api.UnimplementedCourierServer

	store  FeedStore
	parser FeedParser
}

func newRPC(grpcs *grpc.Server, store FeedStore, parser FeedParser) *rpc {
	svc := rpc{store: store, parser: parser}
	api.RegisterCourierServer(grpcs, &svc)
	return &svc
}

// EditFeed satisfies the service API.
func (r *rpc) EditFeed(
	_ context.Context,
	_ *api.EditFeedRequest,
) (*api.EditFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// DeleteFeeds satisfies the service API.
func (r *rpc) DeleteFeeds(
	_ context.Context,
	_ *api.DeleteFeedsRequest,
) (*api.DeleteFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// PollFeeds satisfies the service API.
func (r *rpc) PollFeeds(_ api.Courier_PollFeedsServer) error {
	return status.Errorf(codes.Unimplemented, "unimplemented")
}

// EditEntry satisfies the service API.
func (r *rpc) EditEntry(
	_ context.Context,
	_ *api.EditEntryRequest,
) (*api.EditEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// ExportOPML satisfies the service API.
func (r *rpc) ExportOPML(
	_ context.Context,
	_ *api.ExportOPMLRequest,
) (*api.ExportOPMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

// ImportOPML satisfies the service API.
func (r *rpc) ImportOPML(
	_ context.Context,
	_ *api.ImportOPMLRequest,
) (*api.ImportOPMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
