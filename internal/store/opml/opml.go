// package opml provides functionalities for parsing and writing OPML files.
//
// It follows the OPML 2.0 specifications [1], but keeps only tags relevant to processing
// subscription lists. Elements relating to display settings, such as expansionState or
// vertScrollState, are omitted.
//
// [1] http://opml.org/spec2.opml
package opml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"golang.org/x/net/html/charset"
)

// Parse parses the given raw OPML document into an OPML struct. Only version 2.0 is supported.
func Parse(raw []byte) (*OPML, error) {
	dec := xml.NewDecoder(bytes.NewReader(raw))
	dec.CharsetReader = charset.NewReaderLabel

	var doc OPML
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	if v := doc.Version; v != "2.0" {
		return nil, fmt.Errorf("opml: version '%s' is unsupported", v)
	}

	return &doc, nil
}

// OPML represents the minimal contents of an OPML file required to for storing a subscription list.
type OPML struct {
	Body    Body   `xml:"body"`
	Head    Head   `xml:"head"`
	Version string `xml:"version,attr"`
}

// Head is the <head> element of an OPML file.
type Head struct {
	Title       *string     `xml:"title"`
	DateCreated *rfc822Time `xml:"dateCreated"`
}

// Body is the <body> element of an OPML file.
type Body struct {
	Outlines []*Outline `xml:"outline"`
}

// Outline is a single outline item in the OPML body. It represents a single subscription / feed.
// Nesting is not supported.
type Outline struct {
	Text   string `xml:"text,attr"`
	Type   string `xml:"type,attr"`
	XMLURL string `xml:"xmlUrl,attr"`

	Description *string `xml:"description,attr"`
	HTMLURL     *string `xml:"htmlUrl,attr"`
}

// rfc822Time wraps the time.Time struct for use in UnmarshalXML.
type rfc822Time struct {
	time.Time
}

// tsFormats is an array of possible time formats that can be found in an OPML file. These are
// roughly based on RFC822, with variations in number of digits for day and year, and
// presence/absence of minutes. When parsing, they are iterated over in-order.
var tsFormats = []string{
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 06 15:04:05 MST",
	"Mon, 02 Jan 06 15:04 MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon, 2 Jan 06 15:04 MST",
}

func (t *rfc822Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var raw string
	_ = d.DecodeElement(&raw, &start)

	var (
		ts  time.Time
		err error
	)
	for _, format := range tsFormats {
		ts, err = time.Parse(format, raw)
		if err == nil {
			break
		}
	}
	if ts.IsZero() {
		return fmt.Errorf("opml: invalid time: %s", raw)
	}

	*t = rfc822Time{ts}

	return nil
}
