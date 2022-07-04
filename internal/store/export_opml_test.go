package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bow/courier/internal/store/opml"
)

func TestExportOPMLEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	payload, err := st.ExportOPML(context.Background())
	r.Nil(payload)

	a.ErrorIs(err, opml.ErrEmptyDocument)
}

func TestExportOPMLExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1", IsRead: false},
				{Title: "Entry A2", IsRead: false},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1", IsRead: false},
			},
		},
	}
	st.addFeeds(dbFeeds)

	payload, err := st.ExportOPML(context.Background())
	r.NoError(err)

	a.Regexp(
		regexp.MustCompile(`<\?xml version="1.0" encoding="UTF-8"\?>
<opml version="2.0">
  <head>
    <title>Courier export</title>
    <dateCreated>\d+ [A-Z][a-z]+ \d+ \d+:\d+ .+</dateCreated>
  </head>
  <body>
    <outline text="Feed X" type="rss" xmlUrl="http://x.com/feed.xml"></outline>
    <outline text="Feed A" type="rss" xmlUrl="http://a.com/feed.xml"></outline>
  </body>
</opml>`),
		string(payload),
	)
}
