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
