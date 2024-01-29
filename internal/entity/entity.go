// Copyright (c) 2023-2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import (
	"sort"
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
	return ID(id), nil
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

func fromEntryPbs(pbs []*api.Entry) []*Entry {
	entries := make([]*Entry, 0)
	for _, pb := range pbs {
		if pb == nil {
			continue
		}
		entries = append(entries, FromEntryPb(pb))
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

type compFunc[T any] func(v1, v2 T) int

type sorter[T any] struct {
	items  []T
	compfs []compFunc[T]
}

func ordered[T any]() *sorter[T] {
	return &sorter[T]{compfs: make([]compFunc[T], 0)}
}

func (s *sorter[T]) By(compf ...compFunc[T]) *sorter[T] {
	s.compfs = append(s.compfs, compf...)
	return s
}

func (s *sorter[T]) Len() int {
	return len(s.items)
}

func (s *sorter[T]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *sorter[T]) Less(i, j int) bool {
	p, q := s.items[i], s.items[j]
	var k int
	for k = 0; k < len(s.compfs)-1; k++ {
		comp := s.compfs[k](p, q)
		if comp < 0 {
			return true
		}
		if comp > 0 {
			return false
		}
	}
	return s.compfs[k](p, q) < 0
}

func (s *sorter[T]) Sort(items []T) {
	s.items = items
	sort.Sort(s)
}
