// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"time"

	"github.com/bow/iris/api"
	"github.com/bow/iris/internal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toFeedPb(feed *internal.Feed) *api.Feed {
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

func toFeedPbs(feeds []*internal.Feed) []*api.Feed {
	pbs := make([]*api.Feed, len(feeds))
	for i, feed := range feeds {
		pbs[i] = toFeedPb(feed)
	}
	return pbs
}

func toEntryPb(entry *internal.Entry) *api.Entry {
	return &api.Entry{
		Id:          entry.ID,
		FeedId:      entry.FeedID,
		Title:       entry.Title,
		IsRead:      entry.IsRead,
		ExtId:       entry.ExtID,
		Description: entry.Description,
		Content:     entry.Content,
		Url:         entry.URL,
		PubTime:     toTimestampPb(entry.Published),
		UpdateTime:  toTimestampPb(entry.Updated),
	}
}

func toEntryPbs(entries []*internal.Entry) []*api.Entry {
	pbs := make([]*api.Entry, len(entries))
	for i, entry := range entries {
		pbs[i] = toEntryPb(entry)
	}
	return pbs
}

func toTimestampPb(v *time.Time) *timestamppb.Timestamp {
	if v == nil {
		return nil
	}
	return timestamppb.New(*v)
}
