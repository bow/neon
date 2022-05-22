package atom

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func Parse(raw []byte) (*Feed, error) {
	var doc Feed
	err := xml.Unmarshal(raw, &doc)
	if err != nil {
		return nil, err
	}

	// Remove empty entries ~ necessary since we can not define the empty / zero value for the
	// Entry struct using the XML field tags.
	es, j := make([]*Entry, len(doc.Entries)), 0
	for _, e := range doc.Entries {
		if !e.IsZero() {
			es[j] = e
			j++
		}
	}
	doc.Entries = es[0:j]

	return &doc, nil
}

// Feed follows RFC3287: https://datatracker.ietf.org/doc/html/rfc4287.
type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	XMLBase *string  `xml:"xml:base,attr"`

	Title    Text        `xml:"title"`
	Subtitle *Text       `xml:"subtitle"`
	Links    []*Link     `xml:"link,omitempty"`
	Updated  RFC3399Time `xml:"updated,omitempty"`
	Author   *Person     `xml:"author"`
	Category []*Category `xml:"category"`
	ID       string      `xml:"id"`
	Entries  []*Entry    `xml:"entry,omitempty"`
}

func (f *Feed) GetURI() string {
	for _, link := range f.Links {
		if link.GetRel() == "self" {
			return link.Href
		}
	}
	return ""
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	XMLBase *string  `xml:"xml:base,attr"`

	Title   Text        `xml:"title"`
	Links   []*Link     `xml:"link,omitempty"`
	ID      string      `xml:"id"`
	Updated RFC3399Time `xml:"updated,omitempty"`
	Summary string      `xml:"summary"`
}

func (e *Entry) IsZero() bool {
	return e.Title.Value == "" &&
		len(e.Links) == 0 &&
		e.ID == "" &&
		e.Updated.IsZero() &&
		e.Summary == ""
}

type Person struct {
	XMLBase *string `xml:"xml:base,attr"`

	Name  string  `xml:"name"`
	URI   *string `xml:"uri"`
	Email *string `xml:"email"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	XMLBase *string  `xml:"xml:base,attr"`

	Term   string  `xml:"term"`
	Scheme *string `xml:"scheme"`
	Label  *string `xml:"label"`
}

func (c *Category) IsZero() bool {
	return c.Term == "" &&
		c.Scheme == nil &&
		c.Label == nil
}

type TextType uint8

const (
	PlainText TextType = iota
	HTMLText
	XHTMLText
)

type Text struct {
	Type  TextType
	Value string
}

func (t *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	_ = d.DecodeElement(&raw, &start)

	var typ = PlainText
	for _, attr := range start.Attr {
		if strings.ToLower(attr.Name.Local) != "type" {
			continue
		}
		switch rt := strings.ToLower(attr.Value); rt {
		case "", "text":
			typ = PlainText
		case "html":
			typ = HTMLText
		case "xhtml":
			typ = XHTMLText
		default:
			return fmt.Errorf("invalid 'type' attribute for tag '%s': '%s'", start.Name.Local, rt)
		}
	}

	*t = Text{Value: raw, Type: typ}
	return nil
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	XMLBase *string  `xml:"xml:base,attr"`

	Href     string  `xml:"href,attr"`
	Rel      *string `xml:"rel,attr"`
	Type     *string `xml:"type,attr"`
	Hreflang *string `xml:"hreflang,attr"`
	Title    *string `xml:"title,attr"`
	Length   *int    `xml:"length,attr"`
}

func (l *Link) GetRel() string {
	if l.Rel == nil {
		return "alternate"
	}
	return *l.Rel
}

type RFC3399Time struct {
	time.Time
}

func (t *RFC3399Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	_ = d.DecodeElement(&raw, &start)
	ts, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return err
	}
	*t = RFC3399Time{ts}
	return nil
}
