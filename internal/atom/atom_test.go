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
	a := assert.New(t)

	feed, err := Parse([]byte(raw))
	r.NoError(err)
	//
	a.Equal("John Doe", feed.Author.Name)
	a.Nil(feed.Author.URI)
	a.Nil(feed.Author.Email)
	//
	a.Len(feed.Categories, 0)
	//
	a.Equal("urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6", feed.ID)
	//
	a.Nil(feed.Subtitle)
	//
	a.Equal("Example Feed", feed.Title.Value)
	a.Equal(PlainText, feed.Title.Type)
	//
	a.Equal(2003, feed.Updated.Year())
	a.Equal(time.December, feed.Updated.Month())
	a.Equal(13, feed.Updated.Day())
	a.Equal(18, feed.Updated.Hour())
	a.Equal(30, feed.Updated.Minute())
	a.Equal(2, feed.Updated.Second())
	//
	a.Equal("", feed.GetURI())

	r.Len(feed.Links, 1)
	//
	flink0 := feed.Links[0]
	a.Equal("http://example.org/", flink0.Href)
	a.Nil(flink0.Hreflang)
	a.Nil(flink0.Length)
	a.Nil(flink0.Rel)
	a.Nil(flink0.Title)
	a.Nil(flink0.Type)

	r.Len(feed.Entries, 1)
	//
	entry0 := feed.Entries[0]
	a.Len(entry0.Categories, 0)
	a.Equal("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a", entry0.ID)
	a.Equal(stringp("Some text."), entry0.Summary)
	a.Equal("Atom-Powered Robots Run Amok", entry0.Title.Value)
	a.Equal(PlainText, entry0.Title.Type)
	a.Equal(2003, entry0.Updated.Year())
	a.Equal(time.December, entry0.Updated.Month())
	a.Equal(13, entry0.Updated.Day())
	a.Equal(18, entry0.Updated.Hour())
	a.Equal(30, entry0.Updated.Minute())
	a.Equal(2, entry0.Updated.Second())

	r.Len(entry0.Links, 1)
	//
	elink0 := entry0.Links[0]
	a.Equal("http://example.org/2003/12/13/atom03", elink0.Href)
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
	a := assert.New(t)

	feed, err := Parse([]byte(raw))
	r.NoError(err)
	//
	a.Equal("John Doe", feed.Author.Name)
	a.Nil(feed.Author.URI)
	a.Nil(feed.Author.Email)
	//
	a.Len(feed.Categories, 0)
	//
	a.Equal("urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6", feed.ID)
	//
	a.Nil(feed.Subtitle)
	//
	a.Equal("Example Feed", feed.Title.Value)
	a.Equal(PlainText, feed.Title.Type)
	//
	a.Equal(2003, feed.Updated.Year())
	a.Equal(time.December, feed.Updated.Month())
	a.Equal(13, feed.Updated.Day())
	a.Equal(18, feed.Updated.Hour())
	a.Equal(30, feed.Updated.Minute())
	a.Equal(2, feed.Updated.Second())
	//
	a.Equal("", feed.GetURI())

	r.Len(feed.Links, 0)

	r.Len(feed.Entries, 0)
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
    <category term="misc"/>
    <category term="atom"/>
  </entry>
</feed>
`
	r := require.New(t)
	a := assert.New(t)

	feed, err := Parse([]byte(raw))
	r.NoError(err)
	//
	a.Nil(feed.Author)
	//
	a.Len(feed.Categories, 0)
	//
	a.Equal("tag:example.org,2003:3", feed.ID)
	//
	a.Equal(`
    A <em>lot</em> of effort
    went into making this effortless
  `,
		feed.Subtitle.Value,
	)
	a.Equal(HTMLText, feed.Subtitle.Type)
	//
	a.Equal("dive into mark", feed.Title.Value)
	a.Equal(PlainText, feed.Title.Type)
	//
	a.Equal(2005, feed.Updated.Year())
	a.Equal(time.July, feed.Updated.Month())
	a.Equal(31, feed.Updated.Day())
	a.Equal(12, feed.Updated.Hour())
	a.Equal(29, feed.Updated.Minute())
	a.Equal(29, feed.Updated.Second())
	//
	a.Equal("http://example.org/feed.atom", feed.GetURI())

	r.Len(feed.Links, 2)
	//
	flink0 := feed.Links[0]
	a.Equal("http://example.org/", flink0.Href)
	a.Equal(stringp("en"), flink0.Hreflang)
	a.Equal(stringp("alternate"), flink0.Rel)
	a.Nil(flink0.Title)
	a.Equal(stringp("text/html"), flink0.Type)
	//
	flink1 := feed.Links[1]
	a.Equal("http://example.org/feed.atom", flink1.Href)
	a.Nil(flink1.Hreflang)
	a.Nil(flink1.Length)
	a.Equal(stringp("self"), flink1.Rel)
	a.Nil(flink1.Title)
	a.Equal(stringp("application/atom+xml"), flink1.Type)

	r.Len(feed.Entries, 1)
	entry0 := feed.Entries[0]
	//
	r.Len(entry0.Categories, 2)
	a.Equal("misc", entry0.Categories[0].Term)
	a.Nil(entry0.Categories[0].Label)
	a.Nil(entry0.Categories[0].Scheme)
	a.Equal("atom", entry0.Categories[1].Term)
	a.Nil(entry0.Categories[1].Label)
	a.Nil(entry0.Categories[1].Scheme)
	//
	a.Equal("tag:example.org,2003:3.2397", entry0.ID)
	a.Nil(entry0.Summary)
	a.Equal("Atom draft-07 snapshot", entry0.Title.Value)
	a.Equal(PlainText, entry0.Title.Type)
	a.Equal(2005, entry0.Updated.Year())
	a.Equal(time.July, entry0.Updated.Month())
	a.Equal(31, entry0.Updated.Day())
	a.Equal(12, entry0.Updated.Hour())
	a.Equal(29, entry0.Updated.Minute())
	a.Equal(29, entry0.Updated.Second())

	r.Len(entry0.Links, 2)
	//
	elink0 := entry0.Links[0]
	a.Equal("http://example.org/2005/04/02/atom", elink0.Href)
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
