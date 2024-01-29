// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import (
	"fmt"
	"time"

	"github.com/bow/neon/internal/opml"
)

type Subscription struct {
	Title *string
	Feeds []*Feed
}

func NewSubscriptionFromRawOPML(payload []byte) (*Subscription, error) {

	doc, err := opml.Parse(payload)
	if err != nil {
		return nil, err
	}

	return NewSubscriptionFromOPML(doc)
}

func NewSubscriptionFromOPML(doc *opml.Doc) (*Subscription, error) {

	feeds := make([]*Feed, len(doc.Body.Outlines))
	for i, outl := range doc.Body.Outlines {
		if outl.Text == "" {
			return nil, fmt.Errorf(
				"missing title for feed with URL=%s in OPML document", outl.XMLURL,
			)
		}
		feed := Feed{
			Title:       outl.Text,
			Description: outl.Description,
			FeedURL:     outl.XMLURL,
			SiteURL:     outl.HTMLURL,
			Tags:        outl.Categories,

			ID:         0,
			Subscribed: time.Time{},
			LastPulled: time.Time{},
		}
		if star := outl.IsStarred; star != nil {
			feed.IsStarred = *star
		}
		feeds[i] = &feed
	}

	sub := Subscription{Feeds: feeds}

	return &sub, nil
}

func (sub *Subscription) Export() ([]byte, error) {
	var et = defaultExportTitle
	if sub.Title != nil {
		et = *sub.Title
	}

	doc := opml.New(et, time.Now())
	for _, feed := range sub.Feeds {
		if err := doc.AddOutline(feed); err != nil {
			return nil, err
		}
	}
	return doc.XML()
}
