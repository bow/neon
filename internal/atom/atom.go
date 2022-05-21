package atom

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type TextType uint8

const (
	PlainText TextType = iota
	HTMLText
	XHTMLText
)

func Parse(raw []byte) (*Document, error) {
	var doc Document
	err := xml.Unmarshal(raw, &doc)
	if err != nil {
		return nil, err
	}

	// Remove empty entries ~ necessary since we can't define 'emptiness' for the Entry structs
	// using tags.
	es, j := make([]*Entry, len(doc.Entries)), 0
	for _, e := range doc.Entries {
		if e.IsNotEmpty() {
			es[j] = e
			j++
		}
	}
	doc.Entries = es[0:j]

	return &doc, nil
}

// Document follows RFC3287: https://datatracker.ietf.org/doc/html/rfc4287.
type Document struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`

	Title   *Text        `xml:"title"`
	Links   []*Link      `xml:"link,omitempty"`
	Updated *RFC3399Time `xml:"updated,omitempty"`
	Author  *Author      `xml:"author"`
	ID      string       `xml:"id"`
	Entries []*Entry     `xml:"entry,omitempty"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`

	Title   *Text        `xml:"title"`
	Links   []*Link      `xml:"link,omitempty"`
	ID      string       `xml:"id"`
	Updated *RFC3399Time `xml:"updated,omitempty"`
	Summary string       `xml:"summary"`
}

func (e *Entry) IsNotEmpty() bool {
	return (e.Title != nil && e.Title.Value != "") ||
		len(e.Links) > 0 ||
		e.ID != "" ||
		e.Updated != nil ||
		e.Summary != ""
}

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
	Href    string   `xml:"href,attr"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
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
