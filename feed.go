package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/rcrowley/mergician/html"
	"golang.org/x/net/html/atom"
)

const FeedLength = 10

type Feed struct {
	Author string // author name
	Path   string // feed path within site URL, like "index.atom.xml"
	Title  string // site title
	URL    string // site URL, like "http://example.com"

	Entries []Entry

	mu sync.Mutex
	t  time.Time
}

func (f *Feed) Add(date, path string, n *html.Node) error {
	if filepath.Base(path) == "index.html" {
		path = filepath.Dir(path) + "/"
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	i := sort.Search(len(f.Entries), func(i int) bool { return f.Entries[i].Date < date })
	if i == len(f.Entries) || f.Entries[i].Path != path {
		f.Entries = append(f.Entries, Entry{})
		copy(f.Entries[i+1:], f.Entries[i:])
		f.Entries[i].Date = date
		f.Entries[i].Path = path

		f.Entries[i].Content = html.Find(n, html.IsAtom(atom.Article))
		if f.Entries[i].Content == nil {
			log.Printf("# no <article> in %s, using class=\"feed\"", path)
			f.Entries[i].Content = html.Find(n, html.All(
				html.Not(html.IsAtom(atom.Time)),
				html.HasAttr("class", "feed"),
			))
			if f.Entries[i].Content.DataAtom == atom.Body {
				f.Entries[i].Content.Data = "div"
				f.Entries[i].Content.DataAtom = atom.Div
			}
		}
		if f.Entries[i].Content == nil {
			return fmt.Errorf("no <article> or element with class=\"feed\" in %s", path)
		}

		f.Entries[i].H1 = html.Find(f.Entries[i].Content, html.IsAtom(atom.H1))
		if f.Entries[i].H1 == nil {
			return fmt.Errorf("no <h1> in %s", path)
		}

	}
	return nil
}

func (f *Feed) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	t := time.Now()
	if !f.t.IsZero() {
		t = f.t
	}
	u, err := url.Parse(f.URL)
	if err != nil {
		return err
	}

	e.EncodeToken(xml.Header)
	e.EncodeToken(xml.StartElement{xml.Name{Local: "feed"}, []xml.Attr{{xml.Name{Local: "xmlns"}, "http://www.w3.org/2005/Atom"}}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "author"}, nil})
	e.EncodeToken(xml.StartElement{xml.Name{Local: "name"}, nil})
	e.EncodeToken(xml.CharData(f.Author))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "name"}})
	e.EncodeToken(xml.EndElement{xml.Name{Local: "author"}})

	u.Path = "/" // at the advice of <https://validator.w3.org/feed/>
	e.EncodeToken(xml.StartElement{xml.Name{Local: "id"}, nil})
	e.EncodeToken(xml.CharData(u.String()))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "id"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, []xml.Attr{
		{xml.Name{Local: "href"}, u.String()},
		{xml.Name{Local: "rel"}, "alternate"},
	}})
	e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}}) // encoding/xml doesn't support self-closing tags
	if f.Path != "" {
		u.Path = f.Path
		e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, []xml.Attr{
			{xml.Name{Local: "href"}, u.String()},
			{xml.Name{Local: "rel"}, "self"},
		}})
		e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}}) // encoding/xml doesn't support self-closing tags
	}

	e.EncodeToken(xml.StartElement{xml.Name{Local: "title"}, nil})
	e.EncodeToken(xml.CharData(f.Title))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "title"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "updated"}, nil})
	e.EncodeToken(xml.CharData(t.Format(time.RFC3339)))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "updated"}})

	var i int
	for _, entry := range f.Entries {
		if i == FeedLength {
			break
		}
		i++

		e.EncodeToken(xml.StartElement{xml.Name{Local: "entry"}, nil})

		u.Path = entry.Path
		e.EncodeToken(xml.StartElement{xml.Name{Local: "id"}, nil})
		e.EncodeToken(xml.CharData(u.String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "id"}})
		e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, []xml.Attr{
			{xml.Name{Local: "href"}, u.String()},
			{xml.Name{Local: "rel"}, "alternate"},
		}})
		e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}}) // encoding/xml doesn't support self-closing tags

		e.EncodeToken(xml.StartElement{xml.Name{Local: "title"}, nil})
		e.EncodeToken(xml.CharData(html.Text(entry.H1).String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "title"}})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "updated"}, nil})
		if t, err := time.Parse(time.DateTime, entry.Date); err == nil {
			e.EncodeToken(xml.CharData(t.Format(time.RFC3339)))
		} else if t, err := time.Parse(time.DateOnly, entry.Date); err == nil {
			e.EncodeToken(xml.CharData(t.Format(time.RFC3339)))
		} else {
			log.Printf("error parsing date %q: %v", entry.Date, err)
			e.EncodeToken(xml.CharData(entry.Date))
		}
		e.EncodeToken(xml.EndElement{xml.Name{Local: "updated"}})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "content"}, []xml.Attr{{xml.Name{Local: "type"}, "html"}}})
		e.EncodeToken(xml.CharData(html.String(entry.Content)))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "content"}})

		e.EncodeToken(xml.EndElement{xml.Name{Local: "entry"}})
	}
	e.EncodeToken(xml.EndElement{xml.Name{Local: "feed"}})
	e.EncodeToken(xml.CharData("\n"))

	e.Flush()
	return nil
}

func (f *Feed) Print() {
	must(f.Render(os.Stdout))
}

func (f *Feed) Render(w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "\t")
	return e.Encode(f)
}

func (f *Feed) RenderFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return f.Render(file)
}
