package htmldifflib

import (
	"io"
	"strings"

	"github.com/yosssi/gohtml"

	"golang.org/x/net/html"
)

type Diff struct {
	IsDifferent bool
	Tags1       []string
	Tags2       []string
}

func RenderNode(w io.Writer, n *html.Node) error {
	return html.Render(w, n)
}

func ParseDoc(r io.Reader) (*html.Node, error) {
	return html.Parse(r)
}

func DiffStructure(n1, n2 *html.Node) (Diff, error) {
	var d Diff
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
					c1.Attr = append(c1.Attr, html.Attribute{
						Key: "data-diff",
					})

					if c2 != nil {
						c2.Attr = append(c2.Attr, html.Attribute{
							Key: "data-diff",
						})
					}
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
