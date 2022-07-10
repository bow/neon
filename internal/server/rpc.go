package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bow/courier/api"
	"github.com/bow/courier/internal"
	"github.com/bow/courier/internal/store"
)

// rpc implements the Courier rpc API.
type rpc struct {
	api.UnimplementedCourierServer

	store store.FeedStore
}

func newRPC(grpcs *grpc.Server, str store.FeedStore) *rpc {
	svc := rpc{store: str}
	api.RegisterCourierServer(grpcs, &svc)
	return &svc
}

// AddFeed satisfies the service API.
func (r *rpc) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	created, err := r.store.AddFeed(
		ctx,
		req.GetUrl(),
		req.Title,
		req.Description,
		req.GetTags(),
		req.IsStarred,
	)
	if err != nil {
		return nil, err
	}

	payload, err := created.Proto()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rsp := api.AddFeedResponse{Feed: payload}

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
		proto, err := feed.Proto()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		rsp.Feeds = append(rsp.Feeds, proto)
	}

	return &rsp, nil
}

// EditFeeds satisfies the service API.
func (r *rpc) EditFeeds(
	ctx context.Context,
	req *api.EditFeedsRequest,
) (*api.EditFeedsResponse, error) {

	ops := make([]*store.FeedEditOp, len(req.Ops))
	for i, op := range req.GetOps() {
		ops[i] = store.NewFeedEditOp(op)
	}

	feeds, err := r.store.EditFeeds(ctx, ops)
	if err != nil {
		return nil, err
	}

	rsp := api.EditFeedsResponse{}
	for _, feed := range feeds {
		fp, err := feed.Proto()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		rsp.Feeds = append(rsp.Feeds, fp)
	}

	return &rsp, nil
}

// DeleteFeeds satisfies the service API.
func (r *rpc) DeleteFeeds(
	ctx context.Context,
	req *api.DeleteFeedsRequest,
) (*api.DeleteFeedsResponse, error) {

	ids := make([]store.DBID, len(req.GetFeedIds()))
	for i, id := range req.GetFeedIds() {
		ids[i] = store.DBID(id)
	}

	err := r.store.DeleteFeeds(ctx, ids)

	rsp := api.DeleteFeedsResponse{}

	return &rsp, err
}

// PollFeeds satisfies the service API.
func (r *rpc) PollFeeds(_ api.Courier_PollFeedsServer) error {
	return status.Errorf(codes.Unimplemented, "unimplemented")
}

// EditEntries satisfies the service API.
func (r *rpc) EditEntries(
	ctx context.Context,
	req *api.EditEntriesRequest,
) (*api.EditEntriesResponse, error) {

	ops := make([]*store.EntryEditOp, len(req.Ops))
	for i, op := range req.GetOps() {
		ops[i] = store.NewEntryEditOp(op)
	}

	entries, err := r.store.EditEntries(ctx, ops)
	if err != nil {
		return nil, err
	}

	rsp := api.EditEntriesResponse{}
	for _, entry := range entries {
		ep, err := entry.Proto()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		rsp.Entries = append(rsp.Entries, ep)
	}

	return &rsp, nil
}

// ExportOPML satisfies the service API.
func (r *rpc) ExportOPML(
	ctx context.Context,
	req *api.ExportOPMLRequest,
) (*api.ExportOPMLResponse, error) {

	payload, err := r.store.ExportOPML(ctx, req.Title)
	if err != nil {
		return nil, err
	}

	rsp := api.ExportOPMLResponse{Payload: payload}

	return &rsp, nil
}

// ImportOPML satisfies the service API.
func (r *rpc) ImportOPML(
	ctx context.Context,
	req *api.ImportOPMLRequest,
) (*api.ImportOPMLResponse, error) {

	n, err := r.store.ImportOPML(ctx, req.Payload)
	if err != nil {
		return nil, err
	}

	rsp := api.ImportOPMLResponse{NumImported: int32(n)}

	return &rsp, nil
}

// GetInfo satisfies the service API.
func (r *rpc) GetInfo(
	_ context.Context,
	_ *api.GetInfoRequest,
) (*api.GetInfoResponse, error) {

	rsp := api.GetInfoResponse{
		Name:      internal.AppName(),
		Version:   internal.Version(),
		GitCommit: internal.GitCommit(),
		BuildTime: internal.BuildTime(),
	}

	return &rsp, nil
}
