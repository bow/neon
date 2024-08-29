// Copyright (c) 2023-2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import (
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bow/neon/api"
)

type ID = uint32

func ToFeedID(raw string) (ID, error) {
	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, FeedNotFoundError{ID: raw}
	}
	return ID(id), nil // #nosec: G115
}

func ToFeedIDs(raw []string) ([]ID, error) {
	ids := make([]ID, 0)
	for _, item := range raw {
		id, err := ToFeedID(item)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func FromFeedPb(pb *api.Feed) *Feed {
	if pb == nil {
		return nil
	}
	return &Feed{
		ID:          pb.GetId(),
		Title:       pb.GetTitle(),
		Description: pb.Description,
		FeedURL:     pb.GetFeedUrl(),
		SiteURL:     pb.SiteUrl,
		Subscribed:  *FromTimestampPb(pb.GetSubTime()),
		LastPulled:  *FromTimestampPb(pb.GetLastPullTime()),
		Updated:     FromTimestampPb(pb.GetUpdateTime()),
		IsStarred:   pb.GetIsStarred(),
		Tags:        pb.GetTags(),
		Entries:     fromEntryPbs(pb.GetEntries()),
	}
}

func FromFeedPbs(pbs []*api.Feed) []*Feed {
	feeds := make([]*Feed, 0)
	for _, pb := range pbs {
		if feed := FromFeedPb(pb); feed != nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func FromEntryPb(pb *api.Entry) *Entry {
	if pb == nil {
		return nil
	}
	return &Entry{
		ID:           pb.GetId(),
		FeedID:       pb.GetFeedId(),
		Title:        pb.GetTitle(),
		IsRead:       pb.GetIsRead(),
		IsBookmarked: pb.GetIsBookmarked(),
		ExtID:        pb.GetExtId(),
		Updated:      FromTimestampPb(pb.GetUpdateTime()),
		Published:    FromTimestampPb(pb.GetUpdateTime()),
		Description:  pb.Description,
		Content:      pb.Content,
		URL:          pb.Url,
	}
}

func fromEntryPbs(pbs []*api.Entry) map[ID]*Entry {
	entries := make(map[ID]*Entry)
	for _, pb := range pbs {
		if pb == nil {
			continue
		}
		entry := FromEntryPb(pb)
		entries[entry.ID] = entry
	}
	return entries
}

func FromStatsPb(pb *api.GetStatsResponse_Stats) *Stats {
	return &Stats{
		NumFeeds:             pb.GetNumFeeds(),
		NumEntries:           pb.GetNumEntries(),
		NumEntriesUnread:     pb.GetNumEntriesUnread(),
		LastPullTime:         FromTimestampPb(pb.GetLastPullTime()),
		MostRecentUpdateTime: FromTimestampPb(pb.GetMostRecentUpdateTime()),
	}
}

func FromTimestampPb(pb *timestamppb.Timestamp) *time.Time {
	if pb == nil {
		return nil
	}
	v := pb.AsTime()
	return &v
}

const defaultExportTitle = "neon export"
