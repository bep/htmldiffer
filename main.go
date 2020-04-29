package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bep/htmldiffer/htmldifflib"

	arg "github.com/alexflint/go-arg"
)

type config struct {
	Dir1 string `arg:"required" help:"The first dir"`
	Dir2 string `arg:"required" help:"The second dir"`
}

func main() {
	var cfg config
	arg.MustParse(&cfg)

	if err := process(cfg); err != nil {
		log.Fatal(err)
	}
}

type diff struct {
	path   string
	reason string
}

func process(cfg config) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	out := filepath.Join(dir, "out")
	outDir1, outDir2 := filepath.Join(out, "dir1"), filepath.Join(out, "dir2")
	must(os.RemoveAll(out))
	must(os.MkdirAll(outDir1, 0777))
	must(os.MkdirAll(outDir2, 0777))

	var diffs []diff

	err = filepath.Walk(cfg.Dir1, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f1, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f1.Close()

		rel := strings.TrimPrefix(path, cfg.Dir1)

		filename2 := filepath.Join(cfg.Dir2, rel)
		f2, err := os.Open(filename2)
		if err != nil {
			diffs = append(diffs, diff{path: rel, reason: "Not found"})
			return nil
		}
		defer f2.Close()

		b1, b2 := mustReadAll(f1), mustReadAll(f2)
		b1, b2 = htmldifflib.Format(b1), htmldifflib.Format(b2)

		isDiff, err := htmldifflib.IsDifferent(bytes.NewReader(b1), bytes.NewReader(b2))
		if err != nil {
			return err
		}

		if isDiff {
			diffs = append(diffs, diff{path: rel, reason: "Different"})
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Diffs", diffs)

	return nil

}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustReadAll(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	must(err)
	return b
}
