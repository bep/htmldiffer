package htmldifflib

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDiff(t *testing.T) {
	c := qt.New(t)

	//	dir := "/Users/bep/dev/sites/kubernetes/website"

	//fn1, fn2 := filepath.Join(dir, "public/docs/concepts/cluster-administration/index.html"), filepath.Join(dir, "public_old/docs/concepts/cluster-administration/index.html")
	fn1, fn2 := "../temp/f1.html", "../temp/f2.html"

	f1, f2 := mustOpen(fn1), mustOpen(fn2)
	defer func() {
		f1.Close()
		f2.Close()
	}()
	//  public/docs/concepts/cluster-administration/index.html vs public_old/docs/concepts/cluster-administration/index.html

	_, err := IsDifferent(f1, f2)
	c.Assert(err, qt.IsNil)

	//formatAndWrite(f1, "temp/f1.html")
	//formatAndWrite(f2, "temp/f2.html")

	//c.Assert(got, qt.Equals, `asdfasdf`)
}

func mustOpen(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return f
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func formatAndWrite(r io.ReadSeeker, targetFilename string) {
	r.Seek(0, 0)
	b, _ := ioutil.ReadAll(r)
	formatted := Format(b)
	current, err := os.Getwd()
	must(err)
	must(ioutil.WriteFile(filepath.Join(current, "..", targetFilename), formatted, 0777))
}
