package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var requiredIDs = []string{
	"product.docs-puller",
	"product.nicos-catalog",
	"product.openbook",
	"product.agent-ops",
	"product.nicos-hidden-menubar",
	"product.jobkit",
}

var forbidden = []string{
	"noise", "getnoise.com", "Bayer", "Enhearten", "EduRAIN", "EduRain",
	"SmartSpectra", "smartspectra", "nvault", "pw-harness", "Garrid",
	"farm-game", "runescape-sim", "MemeBattle", "sol-surfer", "idle-time",
	"1,000+", "30,000+", "3,400+", "$50M", "nicostranquist.com",
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "catalog.yaml")); err != nil {
		t.Fatalf("catalog.yaml not found from %s: %v", wd, err)
	}
	return root
}

func loadTestCatalog(t *testing.T) Catalog {
	t.Helper()
	cat, err := loadCatalog(filepath.Join(repoRoot(t), "catalog.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return cat
}

func TestFeaturedOrder(t *testing.T) {
	cat := loadTestCatalog(t)
	if len(cat.Featured) != len(requiredIDs) {
		t.Fatalf("featured count %d, want %d", len(cat.Featured), len(requiredIDs))
	}
	for i, id := range requiredIDs {
		if cat.Featured[i].ID != id {
			t.Errorf("featured[%d]=%s, want %s", i, cat.Featured[i].ID, id)
		}
	}
}

func TestIndexIsASiteNotAGreeting(t *testing.T) {
	html := renderIndex(loadTestCatalog(t))
	if !strings.Contains(html, `<link rel="stylesheet" href="site.css">`) {
		t.Fatal("index must be a real HTML site")
	}
	if strings.Contains(html, "Hi, I'm") || strings.Contains(html, "material-ui") {
		t.Fatal("index still looks like the 2019 leftover")
	}
	if !strings.Contains(html, "public BM25 sample") {
		t.Fatal("docs-puller sample wording missing")
	}
	if !strings.Contains(html, "synthetic-fixture") {
		t.Fatal("jobkit boundary missing")
	}
	for _, p := range loadTestCatalog(t).Featured {
		if !strings.Contains(html, p.URL) || !strings.Contains(html, p.ProofURL) {
			t.Errorf("missing proof for %s", p.ID)
		}
	}
	for _, bad := range forbidden {
		if strings.Contains(html, bad) {
			t.Errorf("forbidden token %q", bad)
		}
	}
}

func TestRenderCheck(t *testing.T) {
	root := t.TempDir()
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "catalog.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "catalog.yaml"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(root, false); err != nil {
		t.Fatal(err)
	}
	if err := run(root, true); err != nil {
		t.Fatal(err)
	}
}
