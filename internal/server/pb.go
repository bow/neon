// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
)

func toFeedPb(feed *entity.Feed) *api.Feed {
	return &api.Feed{
		Id:           feed.ID,
		Title:        feed.Title,
		FeedUrl:      feed.FeedURL,
		SiteUrl:      feed.SiteURL,
		Tags:         feed.Tags,
		Description:  feed.Description,
		IsStarred:    feed.IsStarred,
		SubTime:      timestamppb.New(feed.Subscribed),
		LastPullTime: timestamppb.New(feed.LastPulled),
		UpdateTime:   toTimestampPb(feed.Updated),
		Entries:      toEntryPbs(feed.Entries),
	}
}

func toFeedPbs(feeds []*entity.Feed) []*api.Feed {
	pbs := make([]*api.Feed, len(feeds))
	for i, feed := range feeds {
		pbs[i] = toFeedPb(feed)
	}
	return pbs
}

func fromFeedEditOpPb(pb *api.EditFeedsRequest_Op) *entity.FeedEditOp {
	return &entity.FeedEditOp{
		ID:          pb.Id,
		Title:       pb.Fields.Title,
		Description: pb.Fields.Description,
		Tags:        &pb.Fields.Tags,
		IsStarred:   pb.Fields.IsStarred,
	}
}

func fromFeedEditOpPbs(pbs []*api.EditFeedsRequest_Op) []*entity.FeedEditOp {
	ops := make([]*entity.FeedEditOp, len(pbs))
	for i, pb := range pbs {
		ops[i] = fromFeedEditOpPb(pb)
	}
	return ops
}

func toEntryPb(entry *entity.Entry) *api.Entry {
	return &api.Entry{
		Id:           entry.ID,
		FeedId:       entry.FeedID,
		Title:        entry.Title,
		IsRead:       entry.IsRead,
		IsBookmarked: entry.IsBookmarked,
		ExtId:        entry.ExtID,
		Description:  entry.Description,
		Content:      entry.Content,
		Url:          entry.URL,
		PubTime:      toTimestampPb(entry.Published),
		UpdateTime:   toTimestampPb(entry.Updated),
	}
}

func toEntryPbs(entries []*entity.Entry) []*api.Entry {
	pbs := make([]*api.Entry, len(entries))
	for i, entry := range entries {
		pbs[i] = toEntryPb(entry)
	}
	return pbs
}

func fromEntryEditOpPb(pb *api.EditEntriesRequest_Op) *entity.EntryEditOp {
	return &entity.EntryEditOp{
		ID:           pb.Id,
		IsRead:       pb.Fields.IsRead,
		IsBookmarked: pb.Fields.IsBookmarked,
	}
}

func fromEntryEditOpPbs(pbs []*api.EditEntriesRequest_Op) []*entity.EntryEditOp {
	ops := make([]*entity.EntryEditOp, len(pbs))
	for i, pb := range pbs {
		ops[i] = fromEntryEditOpPb(pb)
	}
	return ops
}

func toStatsPb(stats *entity.Stats) *api.GetStatsResponse_Stats {
	return &api.GetStatsResponse_Stats{
		NumFeeds:             stats.NumFeeds,
		NumEntries:           stats.NumEntries,
		NumEntriesUnread:     stats.NumEntriesUnread,
		LastPullTime:         toTimestampPb(stats.LastPullTime),
		MostRecentUpdateTime: toTimestampPb(stats.MostRecentUpdateTime),
	}
}

func toTimestampPb(v *time.Time) *timestamppb.Timestamp {
	if v == nil {
		return nil
	}
	return timestamppb.New(*v)
}
