package htmldifflib

import (
	"io"

	"github.com/yosssi/gohtml"

	"golang.org/x/net/html"
)

func IsDifferent(doc1, doc2 io.Reader) (bool, error) {

	n1, err := html.Parse(doc1)
	if err != nil {
		return false, err
	}
	n2, err := html.Parse(doc2)
	if err != nil {
		return false, err
	}

	isDifferent := false

	var walk func(n1, n2 *html.Node)
	walk = func(n1, n2 *html.Node) {
		for c1, c2 := n1.FirstChild, n2.FirstChild; c1 != nil; c1, c2 = c1.NextSibling, c2.NextSibling {
			if c2 == nil {
				//fmt.Println("C2 is nil", n1.Data, "=>", c1.Data)
				isDifferent = true
				return
			}
			walk(c1, c2)
		}
	}

	walk(n1, n2)

	if isDifferent {
		return true, nil
	}

	walk(n2, n1)
	return isDifferent, nil
}

func Format(b []byte) []byte {
	return gohtml.FormatBytes(b)
}
