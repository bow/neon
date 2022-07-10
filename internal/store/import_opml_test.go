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

	n, err := st.ImportOPML(context.Background(), []byte{})
	r.Equal(0, n)
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

	n, err := st.ImportOPML(context.Background(), []byte(payload))
	r.NoError(err)

	a.Equal(0, n)
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

	n, err := st.ImportOPML(context.Background(), []byte(payload))
	r.NoError(err)

	a.Equal(1, n)
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
			Title:     "Feed BC",
			FeedURL:   "http://bc.com/feed.xml",
			Updated:   WrapNullString("2022-03-19T16:23:18.600+02:00"),
			IsStarred: false,
			Entries: []*Entry{
				{Title: "Entry BC1", IsRead: false},
				{Title: "Entry BC2", IsRead: true},
			},
		},
		{
			Title:     "Feed D",
			FeedURL:   "http://d.com/feed.xml",
			Updated:   WrapNullString("2022-04-20T16:32:30.760+02:00"),
			IsStarred: true,
			Entries: []*Entry{
				{Title: "Entry D1", IsRead: false},
			},
			Tags: []string{"foo", "baz"},
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
		xmlns:courier="https://github.com/bow/courier"
		courier:isStarred="true"
	></outline>
  </body>
</opml>
`

	r.Equal(2, st.countFeeds())
	a.False(existfA())
	a.False(existfBC())

	n, err := st.ImportOPML(context.Background(), []byte(payload))
	r.NoError(err)

	a.Equal(2, n)
	a.Equal(3, st.countFeeds())
	a.True(existfA())
	a.True(existfBC())
}
