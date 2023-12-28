// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import "context"

// Datastore describes the persistence layer interface.
type Datastore interface {
	AddFeed(
		ctx context.Context,
		feedURL string,
		title *string,
		desc *string,
		tags []string,
		isStarred *bool,
	) (
		feed *Feed,
		added bool,
		err error,
	)

	EditFeeds(
		ctx context.Context,
		ops []*FeedEditOp,
	) (
		feeds []*Feed,
		err error,
	)

	ListFeeds(
		ctx context.Context,
	) (
		feeds []*Feed,
		err error,
	)

	PullFeeds(
		ctx context.Context,
		ids []ID,
	) (
		results <-chan PullResult,
	)

	DeleteFeeds(
		ctx context.Context,
		ids []ID,
	) (
		err error,
	)

	ListEntries(
		ctx context.Context,
		feedIDs []ID,
		isBookmarked *bool,
	) (
		entries []*Entry,
		err error,
	)

	EditEntries(
		ctx context.Context,
		ops []*EntryEditOp,
	) (
		entries []*Entry,
		err error,
	)

	GetEntry(
		ctx context.Context,
		id ID,
	) (
		entry *Entry,
		err error,
	)

	ExportSubscription(
		ctx context.Context,
		title *string,
	) (
		subscription *Subscription,
		err error,
	)

	ImportSubscription(
		ctx context.Context,
		sub *Subscription,
	) (
		processed int,
		imported int,
		err error,
	)

	GetGlobalStats(
		ctx context.Context,
	) (
		stats *Stats,
		err error,
	)
}
