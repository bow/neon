// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportOPMLErrEmptyPayload(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	nproc, nimp, err := st.ImportOPML(context.Background(), []byte{})
	r.Equal(0, nproc)
	r.Equal(0, nimp)
	a.EqualError(err, "payload is empty")

	a.Equal(0, st.countFeeds())
}

func TestImportOPMLOkEmptyOPMLBody(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	payload := `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head>
    <title>mySubscriptions.opml</title>
    <dateCreated>Sat, 18 Jun 2005 12:11:52 GMT</dateCreated>
  </head>
  <body>
  </body>
</opml>
`

	nproc, nimp, err := st.ImportOPML(context.Background(), []byte(payload))
	r.NoError(err)

	a.Equal(0, nproc)
	a.Equal(0, nimp)
	a.Equal(0, st.countFeeds())
}

func TestImportOPMLOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	existf := func() bool {
		return st.rowExists(
			feedExistSQL,
			"Feed A",
			nil,
			"http://a.com/feed.xml",
			nil,
			false,
		)
	}

	r.Equal(0, st.countFeeds())
	a.False(existf())

	payload := `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head>
    <title>mySubscriptions.opml</title>
    <dateCreated>Sat, 18 Jun 2005 12:11:52 GMT</dateCreated>
  </head>
  <body>
    <outline text="Feed A" type="rss" xmlUrl="http://a.com/feed.xml"></outline>
  </body>
</opml>
`

	nproc, nimp, err := st.ImportOPML(context.Background(), []byte(payload))
	r.NoError(err)

	a.Equal(1, nproc)
	a.Equal(1, nimp)
	a.Equal(1, st.countFeeds())
	a.True(existf())
}

func TestImportOPMLOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			title:     "Feed BC",
			feedURL:   "http://bc.com/feed.xml",
			updated:   toNullString("2022-03-19T16:23:18.600+02:00"),
			isStarred: false,
			entries: []*Entry{
				{Title: "Entry BC1", IsRead: false},
				{Title: "Entry BC2", IsRead: true},
			},
		},
		{
			title:     "Feed D",
			feedURL:   "http://d.com/feed.xml",
			updated:   toNullString("2022-04-20T16:32:30.760+02:00"),
			isStarred: true,
			entries: []*Entry{
				{Title: "Entry D1", IsRead: false},
			},
			tags: []string{"foo", "baz"},
		},
	}
	st.addFeeds(dbFeeds)

	existfA := func() bool {
		return st.rowExists(
			feedExistSQL,
			"Feed A",
			"New feed",
			"http://a.com/feed.xml",
			"http://a.com",
			false,
		)
	}
	existfBC := func() bool {
		return st.rowExists(
			feedExistSQL,
			"Feed BC",
			"Updated feed",
			"http://bc.com/feed.xml",
			nil,
			true,
		)
	}

	payload := `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head>
    <title>mySubscriptions.opml</title>
    <dateCreated>Sun, 26 Jun 2005 18:08:31 GMT</dateCreated>
  </head>
  <body>
    <outline
		text="Feed A"
		type="rss"
		xmlUrl="http://a.com/feed.xml"
		htmlUrl="http://a.com"
		description="New feed"
	></outline>
	<outline
		text="Feed BC"
		type="rss"
		xmlUrl="http://bc.com/feed.xml"
		description="Updated feed"
		xmlns:iris="https://github.com/bow/iris"
		iris:isStarred="true"
	></outline>
  </body>
</opml>
`

	r.Equal(2, st.countFeeds())
	a.False(existfA())
	a.False(existfBC())

	nproc, nimp, err := st.ImportOPML(context.Background(), []byte(payload))
	r.NoError(err)

	a.Equal(2, nproc)
	a.Equal(1, nimp)
	a.Equal(3, st.countFeeds())
	a.True(existfA())
	a.True(existfBC())
}
