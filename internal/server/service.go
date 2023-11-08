// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"fmt"

	"github.com/bow/iris/api"
	"github.com/bow/iris/internal"
	"github.com/bow/iris/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// service implements the iris service API.
type service struct {
	api.UnimplementedIrisServer

	store store.FeedStore
}

// AddFeed satisfies the service API.
func (svc *service) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	created, err := svc.store.AddFeed(
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

	rsp := api.AddFeedResponse{Feed: toFeedPb(created)}

	return &rsp, nil
}

// ListFeeds satisfies the service API.
func (svc *service) ListFeeds(
	ctx context.Context,
	_ *api.ListFeedsRequest,
) (*api.ListFeedsResponse, error) {

	feeds, err := svc.store.ListFeeds(ctx)
	if err != nil {
		return nil, err
	}
	rsp := api.ListFeedsResponse{Feeds: toFeedPbs(feeds)}

	return &rsp, nil
}

// EditFeeds satisfies the service API.
func (svc *service) EditFeeds(
	ctx context.Context,
	req *api.EditFeedsRequest,
) (*api.EditFeedsResponse, error) {

	ops := fromFeedEditOpPbs(req.GetOps())
	feeds, err := svc.store.EditFeeds(ctx, ops)
	if err != nil {
		return nil, err
	}

	rsp := api.EditFeedsResponse{Feeds: toFeedPbs(feeds)}

	return &rsp, nil
}

// DeleteFeeds satisfies the service API.
func (svc *service) DeleteFeeds(
	ctx context.Context,
	req *api.DeleteFeedsRequest,
) (*api.DeleteFeedsResponse, error) {

	ids := make([]store.ID, len(req.GetFeedIds()))
	for i, id := range req.GetFeedIds() {
		ids[i] = id
	}

	err := svc.store.DeleteFeeds(ctx, ids)

	rsp := api.DeleteFeedsResponse{}

	return &rsp, err
}

// PullFeeds satisfies the service API.
func (svc *service) PullFeeds(
	req *api.PullFeedsRequest,
	stream api.Iris_PullFeedsServer,
) error {

	convert := func(pr internal.PullResult) (*api.PullFeedsResponse, error) {
		if err := pr.Error(); err != nil {
			url := pr.URL()
			if url == "" {
				return nil, err
			}
			rspErr := err.Error()
			rsp := api.PullFeedsResponse{Url: url, Error: &rspErr}
			return &rsp, nil
		}
		feed := pr.Feed()
		if feed == nil {
			return nil, nil
		}
		rsp := api.PullFeedsResponse{Url: pr.URL(), Feed: toFeedPb(feed)}

		return &rsp, nil
	}

	ids := make([]store.ID, len(req.GetFeedIds()))
	for i, id := range req.GetFeedIds() {
		ids[i] = id
	}

	ch := svc.store.PullFeeds(stream.Context(), ids)

	for pr := range ch {
		payload, err := convert(pr)
		if err != nil {
			return err
		}
		if payload == nil {
			continue
		}
		if err := stream.Send(payload); err != nil {
			return err
		}
	}

	return nil
}

// ListEntries satisfies the service API.
func (svc *service) ListEntries(
	ctx context.Context,
	req *api.ListEntriesRequest,
) (*api.ListEntriesResponse, error) {

	entries, err := svc.store.ListEntries(ctx, req.GetFeedId())
	if err != nil {
		return nil, err
	}

	rsp := api.ListEntriesResponse{Entries: toEntryPbs(entries)}

	return &rsp, nil
}

// EditEntries satisfies the service API.
func (svc *service) EditEntries(
	ctx context.Context,
	req *api.EditEntriesRequest,
) (*api.EditEntriesResponse, error) {

	ops := fromEntryEditOpPbs(req.GetOps())
	entries, err := svc.store.EditEntries(ctx, ops)
	if err != nil {
		return nil, err
	}

	rsp := api.EditEntriesResponse{Entries: toEntryPbs(entries)}

	return &rsp, nil
}

// GetEntry satisfies the service API.
func (svc *service) GetEntry(
	ctx context.Context,
	req *api.GetEntryRequest,
) (*api.GetEntryResponse, error) {

	entry, err := svc.store.GetEntry(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	rsp := api.GetEntryResponse{Entry: toEntryPb(entry)}

	return &rsp, nil
}

// ExportOPML satisfies the service API.
func (svc *service) ExportOPML(
	ctx context.Context,
	req *api.ExportOPMLRequest,
) (*api.ExportOPMLResponse, error) {

	sub, err := svc.store.ExportSubscription(ctx, req.Title)
	if err != nil {
		return nil, err
	}

	payload, err := sub.Export()
	if err != nil {
		msg := fmt.Errorf("failed to convert subscriptions to OPML: %w", err).Error()
		return nil, status.Errorf(codes.Internal, msg)
	}

	rsp := api.ExportOPMLResponse{Payload: payload}

	return &rsp, nil
}

// ImportOPML satisfies the service API.
func (svc *service) ImportOPML(
	ctx context.Context,
	req *api.ImportOPMLRequest,
) (*api.ImportOPMLResponse, error) {

	nproc, nimp, err := svc.store.ImportOPML(ctx, req.Payload)
	if err != nil {
		return nil, err
	}

	rsp := api.ImportOPMLResponse{
		NumProcessed: uint32(nproc),
		NumImported:  uint32(nimp),
	}

	return &rsp, nil
}

// GetStats satisfies the service API.
func (svc *service) GetStats(
	ctx context.Context,
	_ *api.GetStatsRequest,
) (*api.GetStatsResponse, error) {

	gstats, err := svc.store.GetGlobalStats(ctx)
	if err != nil {
		return nil, err
	}

	rsp := api.GetStatsResponse{Global: toStatsPb(gstats)}

	return &rsp, nil
}

// GetInfo satisfies the service API.
func (svc *service) GetInfo(
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
