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
	CommonAttributes
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`

	Authors      []*Person   `xml:"author"`
	Categories   []*Category `xml:"category"`
	Contributors []*Person   `xml:"contributor"`
	Entries      []*Entry    `xml:"entry"`
	Generator    *Generator  `xml:"generator"`
	Icon         *string     `xml:"icon"`
	ID           string      `xml:"id"`
	Links        []*Link     `xml:"link"`
	Logo         *string     `xml:"logo"`
	Rights       *string     `xml:"rights"`
	Subtitle     *Text       `xml:"subtitle"`
	Title        Text        `xml:"title"`
	UpdatedUTC   RFC3399Time `xml:"updated"`
}

func (f *Feed) GetPreferredURI() string {
	for _, link := range f.Links {
		if link.GetRel() == "self" {
			return link.Href
		}
	}
	return ""
}

type Entry struct {
	CommonAttributes
	XMLName xml.Name

	Authors      []*Person    `xml:"author"`
	Categories   []*Category  `xml:"category"`
	Content      *Content     `xml:"content"`
	Contributors []*Person    `xml:"contributor"`
	ID           string       `xml:"id"`
	Links        []*Link      `xml:"link"`
	PublishedUTC *RFC3399Time `xml:"published"`
	Source       *Source      `xml:"source"`
	Summary      *Text        `xml:"summary"`
	Title        Text         `xml:"title"`
	UpdatedUTC   RFC3399Time  `xml:"updated"`
}

func (e *Entry) IsZero() bool {
	return e.Title.IsZero() &&
		len(e.Links) == 0 &&
		e.ID == "" &&
		e.UpdatedUTC.IsZero() &&
		e.Summary == nil
}

type CommonAttributes struct {
	XMLBase *string `xml:"base,attr"`
	XMLLang *string `xml:"lang,attr"`
}

type Category struct {
	CommonAttributes
	XMLName xml.Name

	Label  *string `xml:"label,attr"`
	Scheme *string `xml:"scheme,attr"`
	Term   string  `xml:"term,attr"`
}

func (c *Category) IsZero() bool {
	return c.Term == "" &&
		c.Scheme == nil &&
		c.Label == nil
}

type Content struct {
	CommonAttributes
	XMLName xml.Name

	Src   *string `xml:"src,attr"`
	Type  *string `xml:"type,attr"`
	Value string  `xml:",innerxml"`
}

type Generator struct {
	CommonAttributes
	XMLName xml.Name

	URI     *string `xml:"uri,attr"`
	Version *string `xml:"version,attr"`
	Value   string  `xml:",innerxml"`
}

type Link struct {
	CommonAttributes
	XMLName xml.Name

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

type Source struct {
	CommonAttributes

	Authors      []*Person   `xml:"author"`
	Categories   []*Category `xml:"category"`
	Contributors []*Person   `xml:"contributor"`
	Generator    *Generator  `xml:"generator"`
	Icon         *string     `xml:"icon"`
	ID           *string     `xml:"id"`
	Links        []*Link     `xml:"link"`
	Logo         *string     `xml:"logo"`
	Rights       *string     `xml:"rights"`
	Subtitle     *Text       `xml:"subtitle"`
	Title        Text        `xml:"title"`
	UpdatedUTC   RFC3399Time `xml:"updated"`
}

func (s *Source) IsZero() bool {
	return false
}

type Person struct {
	CommonAttributes

	Email *string `xml:"email"`
	Name  string  `xml:"name"`
	URI   *string `xml:"uri"`
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
	*t = RFC3399Time{ts.UTC()}
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
