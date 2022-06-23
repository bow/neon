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

// AddFeed satisfies the service API.
func (r *rpc) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	feed, err := r.parser.ParseURLWithContext(req.GetUrl(), ctx)
	if err != nil {
		return nil, err
	}

	err = r.store.AddFeed(ctx, feed, req.Title, req.Description, req.GetCategories())
	if err != nil {
		return nil, err
	}

	rsp := api.AddFeedResponse{}

	return &rsp, nil
}

// ListFeeds satisfies the service API.
func (r *rpc) ListFeeds(
	ctx context.Context,
	_ *api.ListFeedsRequest,
) (*api.ListFeedsResponse, error) {

	feeds, err := r.store.ListFeeds(ctx)
	if err != nil {
		return nil, err
	}

	rsp := api.ListFeedsResponse{}
	for _, feed := range feeds {
		// TODO: Use gRPC INTERNAL error for this.
		proto, err := feed.Proto()
		if err != nil {
			return nil, err
		}
		rsp.Feeds = append(rsp.Feeds, proto)
	}

	return &rsp, nil
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

// SetEntryFields satisfies the service API.
func (r *rpc) SetEntryFields(
	_ context.Context,
	_ *api.SetEntryFieldsRequest,
) (*api.SetEntryFieldsResponse, error) {
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

// GetInfo satisfies the service API.
func (r *rpc) GetInfo(
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
