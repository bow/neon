package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditEntriesEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	entries, err := st.EditEntries(context.Background(), nil)
	r.NoError(err)

	a.Empty(entries)
}

func TestEditEntriesMinimal(t *testing.T) {
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
				{Title: "Entry A1", IsRead: true},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(1, st.countFeeds())

	existe := func(title string, isRead bool) bool {
		return st.rowExists(
			`SELECT * FROM entries e WHERE e.title = ? AND e.is_read = ?`,
			title,
			isRead,
		)
	}

	a.True(existe("Entry A1", true))
	a.False(existe("Entry A1", false))

	setOps := []*EntryEditOp{
		{DBID: keys["Feed A"].Entries["Entry A1"], IsRead: pointer(false)},
	}
	entries, err := st.EditEntries(context.Background(), setOps)
	r.NoError(err)

	a.Len(entries, 1)

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))
}

func TestEditEntriesExtended(t *testing.T) {
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
	keys := st.addFeeds(dbFeeds)

	r.Equal(2, st.countFeeds())

	existe := func(title string, isRead bool) bool {
		return st.rowExists(
			`SELECT * FROM entries e WHERE e.title = ? AND e.is_read = ?`,
			title,
			isRead,
		)
	}

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))

	a.True(existe("Entry A2", false))
	a.False(existe("Entry A2", true))

	a.True(existe("Entry X1", false))
	a.False(existe("Entry X1", true))

	setOps := []*EntryEditOp{
		{DBID: keys["Feed X"].Entries["Entry X1"], IsRead: pointer(true)},
		{DBID: keys["Feed A"].Entries["Entry A2"], IsRead: pointer(true)},
	}
	entries, err := st.EditEntries(context.Background(), setOps)
	r.NoError(err)

	a.Len(entries, 2)

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))

	a.False(existe("Entry A2", false))
	a.True(existe("Entry A2", true))

	a.False(existe("Entry X1", false))
	a.True(existe("Entry X1", true))
}
