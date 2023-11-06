// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"time"

	"github.com/bow/iris/api"
	"github.com/bow/iris/internal/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toFeedPB(sf *store.FeedRecord) (af *api.Feed, err error) {
	af = &api.Feed{
		Id:           sf.ID(),
		Title:        sf.Title(),
		FeedUrl:      sf.FeedURL(),
		SiteUrl:      sf.SiteURL(),
		Tags:         sf.Tags(),
		Description:  sf.Description(),
		IsStarred:    sf.IsStarred(),
		SubTime:      timestamppb.New(sf.Subscribed()),
		LastPullTime: timestamppb.New(sf.LastPulled()),
		UpdateTime:   toTimestampPB(sf.Updated()),
	}

	for _, entry := range sf.Entries() {
		ep, err := entry.Proto()
		if err != nil {
			return nil, err
		}
		af.Entries = append(af.Entries, ep)
	}

	return af, nil
}

func toTimestampPB(v *time.Time) *timestamppb.Timestamp {
	if v == nil {
		return nil
	}
	return timestamppb.New(*v)
}
