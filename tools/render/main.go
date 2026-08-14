package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	check := flag.Bool("check", false, "exit 1 if generated files are stale")
	root := flag.String("root", ".", "repository root")
	flag.Parse()
	if err := run(*root, *check); err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		os.Exit(1)
	}
}

func run(root string, check bool) error {
	cat, err := loadCatalog(filepath.Join(root, "catalog.yaml"))
	if err != nil {
		return err
	}
	files := map[string]string{
		filepath.Join(root, "index.html"):         renderIndex(cat),
		filepath.Join(root, "404.html"):           render404(cat),
		filepath.Join(root, "robots.txt"):         renderRobots(cat),
		filepath.Join(root, "assets", "mark.svg"): renderMark(),
	}
	var stale []string
	for path, body := range files {
		if check {
			existing, err := os.ReadFile(path)
			if err != nil {
				stale = append(stale, path+" (missing)")
				continue
			}
			if !bytes.Equal(existing, []byte(body)) {
				stale = append(stale, path)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
	}
	if check && len(stale) > 0 {
		return fmt.Errorf("generated files are stale:\n  %s", strings.Join(stale, "\n  "))
	}
	return nil
}
