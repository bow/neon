package opml

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOkExtended(t *testing.T) {
	raw := `<?xml version="1.0" encoding="ISO-8859-1"?>
<opml version="2.0">
  <head>
    <title>mySubscriptions.opml</title>
    <dateCreated>Sat, 18 Jun 2005 12:11:52 GMT</dateCreated>
    <dateModified>Tue, 2 Aug 2005 21:42:48 GMT</dateModified>
    <ownerName>Dave Winer</ownerName>
    <ownerEmail>dave@scripting.com</ownerEmail>
    <expansionState />
    <vertScrollState>1</vertScrollState>
    <windowTop>61</windowTop>
    <windowLeft>304</windowLeft>
    <windowBottom>562</windowBottom>
    <windowRight>842</windowRight>
    </head>
  <body>
    <outline
        text="CNET News.com"
        description="Tech news and business reports by CNET News.com. Focused on information technology, core topics include computers, hardware, software, networking, and Internet media."
        htmlUrl="http://news.com.com/"
        language="unknown"
        title="CNET News.com"
        type="rss"
        version="RSS2"
        xmlUrl="http://news.com.com/2547-1_3-0-5.xml"
    />
    <outline
        text="NYT &gt; Business"
        description="Find breaking news &amp; business news on Wall Street, media &amp; advertising, international business, banking, interest rates, the stock market, currencies &amp; funds."
        htmlUrl="http://www.nytimes.com/pages/business/index.html?partner=rssnyt"
        language="unknown"
        title="NYT &gt; Business"
        type="rss"
        version="RSS2"
        xmlUrl="http://www.nytimes.com/services/xml/rss/nyt/Business.xml"
    />
    <outline
        text="Wired News"
        description="Technology, and the way we do business, is changing the world we know. Wired News is a technology - and business-oriented news service feeding an intelligent, discerning audience. What role does technology play in the day-to-day living of your life? Wired News tells you. How has evolving technology changed the face of the international business world? Wired News puts you in the picture."
        htmlUrl="http://www.wired.com/"
        language="unknown"
        title="Wired News"
        type="rss"
        version="RSS"
        xmlUrl="http://www.wired.com/news_drop/netcenter/netcenter.rdf"
    />
    <outline
        text="NYT &gt; Technology"
        description=""
        htmlUrl="http://www.nytimes.com/pages/technology/index.html?partner=rssnyt"
        language="unknown"
        title="NYT &gt; Technology"
        type="rss"
        version="RSS2"
        xmlUrl="http://www.nytimes.com/services/xml/rss/nyt/Technology.xml"
    />
  </body>
</opml>
`

	r := require.New(t)
	a := assert.New(t)

	doc, err := Parse([]byte(raw))
	r.NoError(err)

	r.NotNil(doc.Head)
	head := doc.Head
	//
	a.Equal(stringp("mySubscriptions.opml"), head.Title)
	//
	a.Equal(2005, head.DateCreated.Year())
	a.Equal(time.June, head.DateCreated.Month())
	a.Equal(18, head.DateCreated.Day())
	a.Equal(12, head.DateCreated.Hour())
	a.Equal(11, head.DateCreated.Minute())
	a.Equal(52, head.DateCreated.Second())
	//
	a.Equal(2005, head.DateModified.Year())
	a.Equal(time.August, head.DateModified.Month())
	a.Equal(2, head.DateModified.Day())
	a.Equal(21, head.DateModified.Hour())
	a.Equal(42, head.DateModified.Minute())
	a.Equal(48, head.DateModified.Second())

	r.NotNil(doc.Body)
	body := doc.Body
	r.NotNil(body.Outlines)
	outls := body.Outlines
	r.Len(outls, 4)
	//
	outl0 := outls[0]
	a.Equal("CNET News.com", outl0.Text)
	a.Equal("rss", outl0.Type)
	a.Equal("http://news.com.com/2547-1_3-0-5.xml", outl0.XMLURL)
	r.NotNil(outl0.Description)
	a.Equal(
		"Tech news and business reports by CNET News.com."+
			" Focused on information technology, core topics include computers, hardware, software,"+
			" networking, and Internet media.",
		*outl0.Description,
	)
	r.NotNil(outl0.HTMLURL)
	a.Equal("http://news.com.com/", *outl0.HTMLURL)
	//
	outl3 := outls[3]
	a.Equal("NYT > Technology", outl3.Text)
	a.Equal("rss", outl3.Type)
	a.Equal("http://www.nytimes.com/services/xml/rss/nyt/Technology.xml", outl3.XMLURL)
	r.NotNil(outl3.Description)
	a.Equal("", *outl3.Description)
	r.NotNil(outl3.HTMLURL)
	a.Equal("http://www.nytimes.com/pages/technology/index.html?partner=rssnyt", *outl3.HTMLURL)
}

func stringp(value string) *string { return &value }
