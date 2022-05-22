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

	feed, err := Parse([]byte(raw))
	r.NoError(err)

	a := assert.New(t)

	a.Equal("Example Feed", feed.Title.Value)
	a.Equal(PlainText, feed.Title.Type)

	a.Nil(feed.Subtitle)

	a.Len(feed.Category, 0)

	a.Equal(2003, feed.Updated.Year())
	a.Equal(time.December, feed.Updated.Month())
	a.Equal(13, feed.Updated.Day())
	a.Equal(18, feed.Updated.Hour())
	a.Equal(30, feed.Updated.Minute())
	a.Equal(2, feed.Updated.Second())

	a.Equal("John Doe", feed.Author.Name)
	a.Nil(feed.Author.URI)
	a.Nil(feed.Author.Email)

	a.Equal("urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6", feed.ID)

	r.Len(feed.Links, 1)
	link0 := feed.Links[0]
	a.Equal("http://example.org/", link0.Href)
	a.Nil(link0.Rel)
	a.Nil(link0.Type)
	a.Nil(link0.Hreflang)
	a.Nil(link0.Title)
	a.Nil(link0.Length)
	a.Equal("", feed.GetURI())

	r.Len(feed.Entries, 1)
	entry := feed.Entries[0]

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

	feed, err := Parse([]byte(raw))
	r.NoError(err)

	a := assert.New(t)

	a.Equal("Example Feed", feed.Title.Value)
	a.Equal(PlainText, feed.Title.Type)

	a.Nil(feed.Subtitle)

	a.Len(feed.Category, 0)

	a.Equal(2003, feed.Updated.Year())
	a.Equal(time.December, feed.Updated.Month())
	a.Equal(13, feed.Updated.Day())
	a.Equal(18, feed.Updated.Hour())
	a.Equal(30, feed.Updated.Minute())
	a.Equal(2, feed.Updated.Second())

	a.Equal("John Doe", feed.Author.Name)
	a.Nil(feed.Author.URI)
	a.Nil(feed.Author.Email)

	r.Len(feed.Links, 0)
	a.Equal("", feed.GetURI())

	a.Equal("urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6", feed.ID)

	a.Len(feed.Entries, 0)
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
    <link
	  rel="enclosure"
	  type="audio/mpeg"
	  length="1337"
	  href="http://example.org/audio/ph34r_my_podcast.mp3"
	/>
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

	feed, err := Parse([]byte(raw))
	r.NoError(err)

	a := assert.New(t)

	a.Equal("dive into mark", feed.Title.Value)
	a.Equal(PlainText, feed.Title.Type)

	a.Equal(`
    A <em>lot</em> of effort
    went into making this effortless
  `,
		feed.Subtitle.Value)
	a.Equal(HTMLText, feed.Subtitle.Type)

	a.Len(feed.Category, 0)

	a.Equal(2005, feed.Updated.Year())
	a.Equal(time.July, feed.Updated.Month())
	a.Equal(31, feed.Updated.Day())
	a.Equal(12, feed.Updated.Hour())
	a.Equal(29, feed.Updated.Minute())
	a.Equal(29, feed.Updated.Second())

	a.Nil(feed.Author)

	a.Equal("tag:example.org,2003:3", feed.ID)

	r.Len(feed.Links, 2)
	link0 := feed.Links[0]
	a.Equal("http://example.org/", link0.Href)
	a.Equal(stringp("alternate"), link0.Rel)
	a.Equal(stringp("text/html"), link0.Type)
	a.Equal(stringp("en"), link0.Hreflang)
	a.Nil(link0.Title)
	link1 := feed.Links[1]
	a.Equal("http://example.org/feed.atom", link1.Href)
	a.Equal(stringp("self"), link1.Rel)
	a.Equal(stringp("application/atom+xml"), link1.Type)
	a.Nil(link1.Hreflang)
	a.Nil(link1.Title)
	a.Nil(link1.Length)
	a.Equal("http://example.org/feed.atom", feed.GetURI())

	a.Len(feed.Entries, 1)
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
	feed, err := Parse([]byte(raw))

	a := assert.New(t)
	a.Nil(feed)
	a.Error(err)
}

func stringp(value string) *string { return &value }
