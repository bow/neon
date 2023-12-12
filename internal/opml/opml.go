// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

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
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

var (
	ErrEmptyPayload  = errors.New("byte payload is empty")
	ErrEmptyDocument = errors.New("OPML document is empty")
)

// Parse parses the given raw OPML document into an OPML struct. Only version 2.0 is supported.
func Parse(raw []byte) (*Doc, error) {

	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}

	dec := xml.NewDecoder(bytes.NewReader(raw))
	dec.CharsetReader = charset.NewReaderLabel

	var doc Doc
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	if v := doc.Version; v != "2.0" {
		return nil, fmt.Errorf("opml: version '%s' is unsupported", v)
	}

	return &doc, nil
}

// Doc represents the minimal contents of an OPML file required to for storing a subscription list.
type Doc struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

func New(title string, created time.Time) *Doc {
	ts := Timestamp(created)
	t := &title
	if title == "" {
		t = nil
	}
	doc := Doc{
		Version: "2.0",
		Head:    Head{Title: t, DateCreated: &ts},
		Body:    Body{},
	}
	return &doc
}

func (doc *Doc) AddOutline(outl Outliner) error {
	item, err := outl.Outline()
	if err != nil {
		return err
	}
	doc.Body.Outlines = append(doc.Body.Outlines, item)

	return nil
}

func (doc *Doc) Empty() bool {
	return len(doc.Body.Outlines) == 0
}

func (doc *Doc) XML() ([]byte, error) {
	if doc.Empty() {
		return nil, ErrEmptyDocument
	}

	var buf bytes.Buffer
	if _, err := buf.WriteString(xml.Header); err != nil {
		return nil, err
	}

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	if err := enc.Encode(doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Head is the <head> element of an OPML file.
type Head struct {
	Title       *string    `xml:"title"`
	DateCreated *Timestamp `xml:"dateCreated"`
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

	Categories  Categories `xml:"category,attr"`
	Description *string    `xml:"description,attr"`
	HTMLURL     *string    `xml:"htmlUrl,attr"`
	IsStarred   *bool      `xml:"https://github.com/bow/lens isStarred,attr,omitempty"`
}

type Outliner interface {
	Outline() (*Outline, error)
}

type Categories []string

const categorySep = ","

func (c *Categories) UnmarshalXMLAttr(attr xml.Attr) error {

	toks := make([]string, 0)
	for _, rt := range strings.Split(attr.Value, categorySep) {
		tok := strings.TrimSpace(rt)
		if tok != "" {
			toks = append(toks, tok)
		}
	}

	*c = Categories(toks)

	return nil
}

func (c *Categories) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	cats := []string(*c)
	if len(cats) == 0 {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: strings.Join(cats, categorySep)}, nil
}

type Timestamp time.Time

func (t *Timestamp) Time() time.Time { return time.Time(*t) }

func (t *Timestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

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
	if err != nil {
		return fmt.Errorf("opml: invalid time: %q matches no expected formats", raw)
	}
	if ts.IsZero() {
		return fmt.Errorf("opml: invalid time: %q is empty", raw)
	}

	*t = Timestamp(ts)

	return nil
}

func (t *Timestamp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tv := time.Time(*t)
	ts := tv.Format(time.RFC822)
	return e.EncodeElement(ts, start)
}

// tsFormats is an array of possible time formats that can be found in an OPML file. These are
// roughly based on RFC822, with variations in number of digits for day and year, and
// presence/absence of minutes. When parsing, they are iterated over in-order.
var tsFormats = []string{
	"02 Jan 2006 15:04:05 MST",
	"02 Jan 2006 15:04 MST",
	"02 Jan 06 15:04:05 MST",
	"02 Jan 06 15:04 MST",
	"2 Jan 2006 15:04:05 MST",
	"2 Jan 2006 15:04 MST",
	"2 Jan 06 15:04:05 MST",
	"2 Jan 06 15:04 MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 06 15:04:05 MST",
	"Mon, 02 Jan 06 15:04 MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon, 2 Jan 06 15:04 MST",
}
