// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/datastore"
	"github.com/bow/neon/internal/entity"
)

// service implements the service API.
type service struct {
	api.UnimplementedNeonServer

	ds datastore.Datastore
}

// AddFeed satisfies the service API.
func (svc *service) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	record, added, err := svc.ds.AddFeed(
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

	rsp := api.AddFeedResponse{Feed: toFeedPb(record), IsAdded: added}

	return &rsp, nil
}

// ListFeeds satisfies the service API.
func (svc *service) ListFeeds(
	ctx context.Context,
	_ *api.ListFeedsRequest,
) (*api.ListFeedsResponse, error) {

	feeds, err := svc.ds.ListFeeds(ctx)
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
	feeds, err := svc.ds.EditFeeds(ctx, ops)
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

	ids := make([]entity.ID, len(req.GetFeedIds()))
	for i, id := range req.GetFeedIds() {
		ids[i] = id
	}

	err := svc.ds.DeleteFeeds(ctx, ids)

	rsp := api.DeleteFeedsResponse{}

	return &rsp, err
}

// PullFeeds satisfies the service API.
func (svc *service) PullFeeds(
	req *api.PullFeedsRequest,
	stream api.Neon_PullFeedsServer,
) error {

	convert := func(pr entity.PullResult) (*api.PullFeedsResponse, error) {
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

	ids := make([]entity.ID, len(req.GetFeedIds()))
	for i, id := range req.GetFeedIds() {
		ids[i] = id
	}

	ch := svc.ds.PullFeeds(stream.Context(), ids)

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

	entries, err := svc.ds.ListEntries(ctx, req.GetFeedIds(), req.IsBookmarked)
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
	entries, err := svc.ds.EditEntries(ctx, ops)
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

	entry, err := svc.ds.GetEntry(ctx, req.GetId())
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

	sub, err := svc.ds.ExportSubscription(ctx, req.Title)
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

	payload := req.GetPayload()

	sub, err := entity.NewSubscriptionFromRawOPML(payload)
	if err != nil {
		msg := fmt.Errorf("failed to parse OPML: %w", err).Error()
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	nproc, nimp, err := svc.ds.ImportSubscription(ctx, sub)
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

	gstats, err := svc.ds.GetGlobalStats(ctx)
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
