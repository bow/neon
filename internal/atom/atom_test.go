package atom

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOkSimple(t *testing.T) {
	raw := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">

  <title>Example Feed</title>
  <link href="http://example.org/"/>
  <updated>2003-12-13T18:30:02Z</updated>
  <author>
    <name>John Doe</name>
  </author>
  <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>

  <entry>
    <title>Atom-Powered Robots Run Amok</title>
    <link href="http://example.org/2003/12/13/atom03"/>
    <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
    <updated>2003-12-13T18:30:02Z</updated>
    <summary>Some text.</summary>
  </entry>

</feed>
`
	r := require.New(t)

	doc, err := Parse([]byte(raw))
	r.NoError(err)

	a := assert.New(t)
	a.Equal("Example Feed", doc.Title.Value)
	a.Equal(PlainText, doc.Title.Type)
	a.Equal(2003, doc.Updated.Year())
	a.Equal(time.December, doc.Updated.Month())
	a.Equal(13, doc.Updated.Day())
	a.Equal(18, doc.Updated.Hour())
	a.Equal(30, doc.Updated.Minute())
	a.Equal(2, doc.Updated.Second())
	a.Equal("John Doe", doc.Author.Name)
	a.Equal("urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6", doc.ID)

	r.Len(doc.Links, 1)
	a.Equal("http://example.org/", doc.Links[0].Href)

	r.Len(doc.Entries, 1)
	entry := doc.Entries[0]

	a.Equal("Atom-Powered Robots Run Amok", entry.Title.Value)
	a.Equal(PlainText, entry.Title.Type)
	a.Equal("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a", entry.ID)
	a.Equal(2003, entry.Updated.Year())
	a.Equal(time.December, entry.Updated.Month())
	a.Equal(13, entry.Updated.Day())
	a.Equal(18, entry.Updated.Hour())
	a.Equal(30, entry.Updated.Minute())
	a.Equal(2, entry.Updated.Second())
	a.Equal("Some text.", entry.Summary)

	r.Len(entry.Links, 1)
	a.Equal("http://example.org/2003/12/13/atom03", entry.Links[0].Href)
}

func TestParseOkMinimal(t *testing.T) {
	raw := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">

  <title>Example Feed</title>
  <updated>2003-12-13T18:30:02Z</updated>
  <author>
    <name>John Doe</name>
  </author>
  <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>

  <entry></entry>

</feed>
`
	r := require.New(t)

	doc, err := Parse([]byte(raw))
	r.NoError(err)

	a := assert.New(t)
	a.Equal("Example Feed", doc.Title.Value)
	a.Equal(PlainText, doc.Title.Type)
	a.Equal(2003, doc.Updated.Year())
	a.Equal(time.December, doc.Updated.Month())
	a.Equal(13, doc.Updated.Day())
	a.Equal(18, doc.Updated.Hour())
	a.Equal(30, doc.Updated.Minute())
	a.Equal(2, doc.Updated.Second())
	a.Equal("John Doe", doc.Author.Name)
	a.Equal("urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6", doc.ID)

	a.Len(doc.Entries, 0)
}

func TestParseErrInvalidTime(t *testing.T) {
	raw := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">

  <title>Example Feed</title>
  <updated>2003-12-13T18:30:02</updated>

  <entry>
  </entry>

</feed>
`
	doc, err := Parse([]byte(raw))

	a := assert.New(t)
	a.Nil(doc)
	a.Error(err)
}
