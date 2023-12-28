// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

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
		feed *entity.Feed,
		added bool,
		err error,
	)

	EditFeeds(
		ctx context.Context,
		ops []*entity.FeedEditOp,
	) (
		feeds []*entity.Feed,
		err error,
	)

	ListFeeds(
		ctx context.Context,
	) (
		feeds []*entity.Feed,
		err error,
	)

	PullFeeds(
		ctx context.Context,
		ids []entity.ID,
	) (
		results <-chan entity.PullResult,
	)

	DeleteFeeds(
		ctx context.Context,
		ids []entity.ID,
	) (
		err error,
	)

	ListEntries(
		ctx context.Context,
		feedIDs []entity.ID,
		isBookmarked *bool,
	) (
		entries []*entity.Entry,
		err error,
	)

	EditEntries(
		ctx context.Context,
		ops []*entity.EntryEditOp,
	) (
		entries []*entity.Entry,
		err error,
	)

	GetEntry(
		ctx context.Context,
		id entity.ID,
	) (
		entry *entity.Entry,
		err error,
	)

	ExportSubscription(
		ctx context.Context,
		title *string,
	) (
		subscription *entity.Subscription,
		err error,
	)

	ImportSubscription(
		ctx context.Context,
		sub *entity.Subscription,
	) (
		processed int,
		imported int,
		err error,
	)

	GetGlobalStats(
		ctx context.Context,
	) (
		stats *entity.Stats,
		err error,
	)
}
