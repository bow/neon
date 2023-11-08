// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bow/iris/internal/opml"
)

func TestExportOPMLOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	payload, err := st.ExportOPML(context.Background(), nil)
	r.Nil(payload)

	a.ErrorIs(err, opml.ErrEmptyDocument)
}

func TestExportOPMLOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*feedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			entries: []*entryRecord{
				{title: "Entry A1", isRead: false},
				{title: "Entry A2", isRead: false},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1", isRead: false},
			},
			tags: []string{"foo", "baz"},
		},
		{
			title:     "Feed Q",
			feedURL:   "http://q.com/feed.xml",
			updated:   toNullTime(mustTime(t, "2022-05-02T11:47:33.683+02:00")),
			isStarred: true,
			entries: []*entryRecord{
				{title: "Entry Q1", isRead: false},
			},
		},
	}
	st.addFeeds(dbFeeds)
	r.Equal(3, st.countFeeds())

	payload, err := st.ExportOPML(context.Background(), pointer("Test Export"))
	r.NoError(err)

	a.Regexp(
		regexp.MustCompile(`<\?xml version="1.0" encoding="UTF-8"\?>
<opml version="2.0">
  <head>
    <title>Test Export</title>
    <dateCreated>\d+ [A-Z][a-z]+ \d+ \d+:\d+ .+</dateCreated>
  </head>
  <body>
    <outline text="Feed Q" type="rss" xmlUrl="http://q.com/feed.xml" xmlns:iris="https://github.com/bow/iris" iris:isStarred="true"></outline>
    <outline text="Feed X" type="rss" xmlUrl="http://x.com/feed.xml" category="foo,baz"></outline>
    <outline text="Feed A" type="rss" xmlUrl="http://a.com/feed.xml"></outline>
  </body>
</opml>`),
		string(payload),
	)
}
