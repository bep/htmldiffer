package main

import (
	"bytes"
	"encoding/json"
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
	Dir1   string `arg:"required" help:"The first dir"`
	Dir2   string `arg:"required" help:"The second dir"`
	OutDir string `arg:"required" help:"The output dir"`
}

func main() {
	var cfg config
	arg.MustParse(&cfg)

	if err := process(cfg); err != nil {
		log.Fatal(err)
	}
}

type diff struct {
	Path   string
	Reason string
}

func process(cfg config) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	out := cfg.OutDir
	staticDir := filepath.Join(out, "static")
	layoutsDir := filepath.Join(out, "layouts")
	assetsDir := filepath.Join(out, "assets")

	outStaticDir1, outStaticDir2 := filepath.Join(staticDir, "dir1"), filepath.Join(staticDir, "dir2")
	outAssetsDir1, outAssetsDir2 := filepath.Join(assetsDir, "dir1"), filepath.Join(assetsDir, "dir2")

	//must(os.RemoveAll(out))
	must(os.MkdirAll(outStaticDir1, 0777))
	must(os.MkdirAll(outStaticDir2, 0777))
	must(os.MkdirAll(layoutsDir, 0777))

	var diffs []diff
	var processCounter int

	copyToStatic := func(fromDir, toDir, path string) {
		b1, err := ioutil.ReadFile(filepath.Join(fromDir, path))
		if err == nil {
			outFilename := filepath.Join(toDir, path)
			must(os.MkdirAll(filepath.Dir(outFilename), 0777))
			must(ioutil.WriteFile(outFilename, b1, 0777))
		}

	}

	walkStatic := func(fromDir, toDir string) error {
		return filepath.Walk(fromDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel := strings.TrimPrefix(path, fromDir)

			if !strings.HasSuffix(path, "html") {
				// CSS etc. just copy it to get the styling correct.
				copyToStatic(fromDir, toDir, rel)

			}
			return nil
		})
	}

	// Copy the static files as-is for preview
	if err := walkStatic(cfg.Dir1, outStaticDir1); err != nil {
		return err
	}
	if err := walkStatic(cfg.Dir2, outStaticDir2); err != nil {
		return err
	}

	err = filepath.Walk(cfg.Dir1, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, "html") {
			return nil

		}

		rel := strings.TrimPrefix(path, cfg.Dir1)

		if strings.Contains(path, "blog") {
			// For now
			return nil
		}

		processCounter++

		f1, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f1.Close()

		b1 := mustReadAll(f1)
		b1 = htmldifflib.Format(b1)

		filename2 := filepath.Join(cfg.Dir2, rel)
		f2, err := os.Open(filename2)
		if err != nil {
			diffs = append(diffs, diff{Path: rel, Reason: "Not found"})
			mustWrite(b1, filepath.Join(outStaticDir1, rel), filepath.Join(outAssetsDir1, rel))
			return nil
		}
		defer f2.Close()

		b2 := mustReadAll(f2)
		b2 = htmldifflib.Format(b2)

		isDiff, err := htmldifflib.IsDifferent(bytes.NewReader(b1), bytes.NewReader(b2))
		if err != nil {
			return err
		}

		if isDiff {
			diffs = append(diffs, diff{Path: rel, Reason: "Different"})
			mustWrite(b1, filepath.Join(outStaticDir1, rel), filepath.Join(outAssetsDir1, rel))
			mustWrite(b2, filepath.Join(outStaticDir2, rel), filepath.Join(outAssetsDir2, rel))
		}

		return nil
	})

	if err != nil {
		return err
	}

	diffsFilename := filepath.Join(assetsDir, "diffs.json")
	must(err)
	b, err := json.MarshalIndent(diffs, "", "  ")
	must(err)
	must(ioutil.WriteFile(diffsFilename, b, 0777))

	fmt.Printf("Diffs: %d/%d\n", len(diffs), processCounter)

	must(copyFile(filepath.Join(dir, "templates", "index.html"), filepath.Join(layoutsDir, "index.html")))
	writeHugoConfig(out)

	return nil

}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustWrite(b []byte, filenames ...string) {
	for _, filename := range filenames {
		must(os.MkdirAll(filepath.Dir(filename), 0777))
		must(ioutil.WriteFile(filename, b, 0777))
	}
}

func mustReadAll(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	must(err)
	return b
}

func writeHugoConfig(dir string) {
	config := `
baseURL = "https://example.org"
disableKinds = [ "page", "section", "taxonomy", "taxonomyTerm"]

`

	must(ioutil.WriteFile(filepath.Join(dir, "config.toml"), []byte(config), 0777))
}

func copyFile(from, to string) error {
	sf, err := os.Open(from)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(to)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(from)
		if err != nil {
			err = os.Chmod(to, si.Mode())

			if err != nil {
				return err
			}
		}

	}
	return nil
}
