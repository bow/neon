package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditFeedsEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	feeds, err := st.EditFeeds(context.Background(), nil)
	r.NoError(err)

	a.Empty(feeds)
}

func TestEditFeedsMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(1, st.countFeeds())

	existf := func(title string) bool {
		return st.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.False(existf("Feed X"))

	ops := []*FeedEditOp{
		{DBID: keys["Feed A"].DBID, Title: pointer("Feed X")},
	}
	feeds, err := st.EditFeeds(context.Background(), ops)
	r.NoError(err)

	a.Len(feeds, 1)

	a.False(existf("Feed A"))
	a.True(existf("Feed X"))
}
