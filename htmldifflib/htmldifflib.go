package htmldifflib

import (
	"io"
	"strings"

	htmldiff "github.com/documize/html-diff"
	"github.com/yosssi/gohtml"

	"golang.org/x/net/html"
)

type Diff struct {
	IsDifferent bool
	Tags1       []string
	Tags2       []string
}

func DiffStructure(doc1, doc2 io.Reader) (Diff, error) {
	var d Diff

	n1, err := html.Parse(doc1)
	if err != nil {
		return d, err
	}
	n2, err := html.Parse(doc2)
	if err != nil {
		return d, err
	}

	isDifferent := false

	var tags1 []string
	var tags2 []string

	var walk func(n1, n2 *html.Node) bool
	walk = func(n1, n2 *html.Node) bool {
		for c1, c2 := n1.FirstChild, n2.FirstChild; c1 != nil; c1, c2 = c1.NextSibling, c2.NextSibling {
			if c1.Type == html.ElementNode {
				tags1 = append(tags1, c1.Data)
				if c2 != nil {
					tags2 = append(tags2, c2.Data)
				}
				if c2 == nil || c2.Data != c1.Data {
					isDifferent = true
					return true
				}
			}

			if c2 == nil {
				break
			}

			if walk(c1, c2) {
				return true
			}

		}

		return false
	}

	walk(n1, n2)

	if !isDifferent {
		tags2 = tags1[:0]
		tags1 = tags2[:0]
		walk(n2, n1)
	}

	if isDifferent {
		d.Tags1 = tags1
		d.Tags2 = tags2
		d.IsDifferent = true

		return d, nil
	}

	return d, nil
}

func DiffMerge(html1, html2 string) ([]string, error) {
	var cfg = &htmldiff.Config{
		Granularity:  5,
		InsertedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: palegreen;"}},
		DeletedSpan:  []htmldiff.Attribute{{Key: "style", Val: "background-color: lightpink;"}},
		ReplacedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: lightskyblue;"}},
		CleanTags:    []string{""},
	}

	res, err := cfg.HTMLdiff([]string{html1, html2})
	if err != nil {
		return []string{""}, nil
	}
	return res, err
}

var htmlCleaner = strings.NewReplacer(
	"\n</code></pre>", "</code></pre>",
	`<nav id="TableOfContents"></nav>`, "",
)

func Format(s string) string {
	s = htmlCleaner.Replace(s)
	return gohtml.Format(s)
}

func FormatBytes(b []byte) []byte {

	return []byte(Format(string(b)))
}
