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
	removeEmptyItems(&doc.Entries)
	removeEmptyItems(&doc.Categories)
	for _, entry := range doc.Entries {
		removeEmptyItems(&entry.Categories)
	}

	return &doc, nil
}

// Feed follows RFC3287: https://datatracker.ietf.org/doc/html/rfc4287.
type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	XMLBase *string  `xml:"base,attr"`

	Author       *Person     `xml:"author"`
	Categories   []*Category `xml:"category"`
	Contributors []*Person   `xml:"contributor"`
	Entries      []*Entry    `xml:"entry,omitempty"`
	ID           string      `xml:"id"`
	Links        []*Link     `xml:"link,omitempty"`
	Subtitle     *Text       `xml:"subtitle"`
	Title        Text        `xml:"title"`
	Updated      RFC3399Time `xml:"updated,omitempty"`
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
	XMLBase *string  `xml:"base,attr"`

	Author       *Person     `xml:"author"`
	Categories   []*Category `xml:"category"`
	Contributors []*Person   `xml:"contributor"`
	ID           string      `xml:"id"`
	Links        []*Link     `xml:"link,omitempty"`
	Summary      *string     `xml:"summary"`
	Title        Text        `xml:"title"`
	Updated      RFC3399Time `xml:"updated,omitempty"`
}

func (e *Entry) IsZero() bool {
	return e.Title.IsZero() &&
		len(e.Links) == 0 &&
		e.ID == "" &&
		e.Updated.IsZero() &&
		e.Summary == nil
}

type Person struct {
	XMLBase *string `xml:"base,attr"`

	Email *string `xml:"email"`
	Name  string  `xml:"name"`
	URI   *string `xml:"uri"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	XMLBase *string  `xml:"base,attr"`

	Label  *string `xml:"label,attr"`
	Scheme *string `xml:"scheme,attr"`
	Term   string  `xml:"term,attr"`
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

func (t *Text) IsZero() bool {
	return t.Value == ""
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
	XMLBase *string  `xml:"base,attr"`

	Href     string  `xml:"href,attr"`
	Hreflang *string `xml:"hreflang,attr"`
	Length   *int    `xml:"length,attr"`
	Rel      *string `xml:"rel,attr"`
	Title    *string `xml:"title,attr"`
	Type     *string `xml:"type,attr"`
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

type zeroer interface {
	IsZero() bool
}

func removeEmptyItems[T zeroer](arr *[]T) {
	var (
		deref = *arr
		n     = len(deref)
		items = make([]T, n)
		j     = 0
	)
	for i := 0; i < n; i++ {
		if item := deref[i]; !item.IsZero() {
			items[j] = item
			j++
		}
	}
	*arr = items[0:j]
}
