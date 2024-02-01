// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/bow/neon/internal/entity"
	"github.com/bow/neon/internal/sliceutil"
)

type feedStore struct {
	items map[entity.ID]*entity.Feed
}

func newFeedStore() *feedStore {
	lfs := feedStore{items: make(map[entity.ID]*entity.Feed)}
	return &lfs
}

func (lfs *feedStore) feedsByPeriod() []feedGroup[feedUpdatePeriod] {
	m := make(map[feedUpdatePeriod][]*entity.Feed)
	for _, feed := range lfs.items {
		key := whenUpdated(feed)
		m[key] = append(m[key], feed)
	}

	groups := make([]feedGroup[feedUpdatePeriod], 0)
	for i := uint8(0); i < uint8(updatedUnknown); i++ {
		period := feedUpdatePeriod(i)
		feeds, hasFeeds := m[period]
		if !hasFeeds {
			continue
		}
		groups = append(groups, newFeedGroup(period, feeds))
	}

	return groups
}

func (lfs *feedStore) upsert(incoming *entity.Feed) {
	if incoming == nil {
		return
	}

	existing, exists := lfs.items[incoming.ID]
	if !exists {
		lfs.items[incoming.ID] = incoming
		return
	}
	lfs.merge(existing, incoming)
}

func (lfs *feedStore) merge(existing, incoming *entity.Feed) {
	existing.Title = incoming.Title
	existing.Description = incoming.Description
	existing.FeedURL = incoming.FeedURL
	existing.SiteURL = incoming.SiteURL
	existing.Subscribed = incoming.Subscribed
	existing.LastPulled = incoming.LastPulled
	existing.Updated = incoming.Updated
	existing.IsStarred = incoming.IsStarred
	existing.Tags = incoming.Tags

	for eid, e := range incoming.Entries {
		existing.Entries[eid] = e
	}
}

type feedGroup[T any] struct {
	label T
	items []*entity.Feed
}

func newFeedGroup[T any](label T, items []*entity.Feed) feedGroup[T] {
	return feedGroup[T]{label: label, items: items}
}

func (fg feedGroup[T]) feedsSlice() []*entity.Feed {
	hasUnread := func(f1, f2 *entity.Feed) int {
		n1, n2 := f1.NumEntriesUnread(), f2.NumEntriesUnread()
		if n1 > 0 && n2 <= 0 {
			return -1
		}
		if n2 > 0 && n1 <= 0 {
			return 1
		}
		return 0
	}
	updateTime := func(f1, f2 *entity.Feed) int {
		ut1, ut2 := f1.Updated, f2.Updated
		if ut1 != nil && ut2 != nil {
			if ut1.Before(*ut2) {
				return 1
			}
			if ut2.Before(*ut1) {
				return -1
			}
			return 0
		}
		if ut1 != nil {
			return -1
		}
		if ut2 != nil {
			return 1
		}
		return 0
	}

	sliceutil.Ordered[*entity.Feed]().
		By(hasUnread, updateTime).
		Sort(fg.items)

	return fg.items
}
