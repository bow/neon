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

	a.Nil(doc.Subtitle)

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

func TestParseOkExtended(t *testing.T) {
	raw := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title type="text">dive into mark</title>
  <subtitle type="html">
    A &lt;em&gt;lot&lt;/em&gt; of effort
    went into making this effortless
  </subtitle>
  <updated>2005-07-31T12:29:29Z</updated>
  <id>tag:example.org,2003:3</id>
  <link rel="alternate" type="text/html" hreflang="en" href="http://example.org/"/>
  <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
  <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  <generator uri="http://www.example.com/" version="1.0">
    Example Toolkit
  </generator>
  <entry>
    <title>Atom draft-07 snapshot</title>
    <link rel="alternate" type="text/html" href="http://example.org/2005/04/02/atom"/>
    <link rel="enclosure" type="audio/mpeg" length="1337" href="http://example.org/audio/ph34r_my_podcast.mp3"/>
    <id>tag:example.org,2003:3.2397</id>
    <updated>2005-07-31T12:29:29Z</updated>
    <published>2003-12-13T08:29:29-04:00</published>
    <author>
      <name>Mark Pilgrim</name>
      <uri>http://example.org/</uri>
      <email>f8dy@example.com</email>
    </author>
    <contributor>
      <name>Sam Ruby</name>
    </contributor>
    <contributor>
      <name>Joe Gregorio</name>
    </contributor>
    <content type="xhtml" xml:lang="en" xml:base="http://diveintomark.org/">
      <div xmlns="http://www.w3.org/1999/xhtml">
        <p><i>[Update: The Atom draft is finished.]</i></p>
      </div>
    </content>
  </entry>
</feed>
`
	r := require.New(t)

	doc, err := Parse([]byte(raw))
	r.NoError(err)

	a := assert.New(t)

	a.Equal("dive into mark", doc.Title.Value)
	a.Equal(PlainText, doc.Title.Type)

	a.Equal(`
    A <em>lot</em> of effort
    went into making this effortless
  `,
		doc.Subtitle.Value)
	a.Equal(HTMLText, doc.Subtitle.Type)

	a.Len(doc.Entries, 1)
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
