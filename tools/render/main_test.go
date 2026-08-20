package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var requiredIDs = []string{
	"product.docs-puller",
	"product.nicos-catalog",
	"product.openbook",
	"product.agent-ops",
	"product.nicos-hidden-menubar",
	"product.nicos-slot-dock",
	"product.jobkit",
	"product.session-pressure",
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
	if !strings.Contains(html, "docs-puller/releases/tag/v0.7.6") {
		t.Fatal("docs-puller release must match the live formal tag")
	}
	if !strings.Contains(html, "https://docs-puller-demo.nstranquist.workers.dev") {
		t.Fatal("docs-puller live demo must be linked")
	}
	if !strings.Contains(html, "session-pressure/releases/tag/v0.1.0") {
		t.Fatal("selected work must include SessionPressure v0.1.0")
	}
	if !strings.Contains(html, "keepawake/releases/tag/v0.1.3") {
		t.Fatal("Also on GitHub must include keepawake v0.1.3")
	}
	if !strings.Contains(html, "nicos-slot-dock/releases/tag/v0.3.6") {
		t.Fatal("selected work must include Nicos Slot Dock v0.3.6")
	}
	if !strings.Contains(html, "id=\"also-public\"") {
		t.Fatal("index must include the Also on GitHub section")
	}
	for _, needle := range []string{"wip-commit", "snapref", "nstranquist/ngtm", "nicos-flag-eval", "nicos-window-switcher"} {
		if !strings.Contains(html, needle) {
			t.Errorf("Also public missing %s", needle)
		}
	}
	if strings.Contains(html, "hearthlight") {
		t.Fatal("hearthlight is not on this public catalog")
	}
	if !strings.Contains(html, "example data only") {
		t.Fatal("jobkit example-data boundary missing")
	}
	if strings.Contains(html, `id="glossary"`) {
		t.Fatal("index must not render a standalone glossary section")
	}
	if !strings.Contains(html, "BM25") {
		t.Fatal("index must define BM25")
	}
	for _, phrase := range []string{
		"inspectable extract",
		"claim boundary",
		"public extracts",
		"AI infrastructure",
		"public extract",
		"stand behind",
		"proof, not posture",
		"unproven claims",
		"fail-closed",
		"missing tags stay visible",
	} {
		if strings.Contains(strings.ToLower(html), phrase) {
			t.Errorf("index still uses factory phrasing %q", phrase)
		}
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

func TestHrefSplitsOutboundAndSameSite(t *testing.T) {
	out := href("https://github.com/nstranquist/docs-puller", "proof")
	if !strings.Contains(out, `target="_blank"`) || !strings.Contains(out, `rel="noopener noreferrer"`) {
		t.Fatalf("outbound href missing new-tab contract: %s", out)
	}
	if !strings.Contains(out, `class="proof"`) {
		t.Fatalf("class not forwarded: %s", out)
	}
	same := href("#work", "")
	if strings.Contains(same, `target="_blank"`) {
		t.Fatalf("in-page hash must stay same-tab: %s", same)
	}
	home := href("./", "brand")
	if strings.Contains(home, `target="_blank"`) {
		t.Fatalf("home link must stay same-tab: %s", home)
	}
}

func TestRenderedPagesLinkContract(t *testing.T) {
	cat := loadTestCatalog(t)
	home := renderIndex(cat)
	missing := render404(cat)
	assertLinkContract(t, "index", home)
	assertLinkContract(t, "404", missing)
	for _, p := range append(append([]Product{}, cat.Featured...), cat.AlsoPublic...) {
		if !strings.Contains(home, `href="`+p.URL+`" target="_blank" rel="noopener noreferrer"`) {
			t.Errorf("index missing new-tab product link for %s", p.ID)
		}
		if !strings.Contains(home, `href="`+p.ProofURL+`" target="_blank" rel="noopener noreferrer"`) &&
			!strings.Contains(home, `href="`+p.ProofURL+`" class="proof" target="_blank" rel="noopener noreferrer"`) {
			t.Errorf("index missing new-tab proof link for %s", p.ID)
		}
	}
	if strings.Contains(missing, `href="#catalog"`) || strings.Contains(missing, `href="#work"`) || strings.Contains(missing, `href="#approach"`) || strings.Contains(missing, `href="#glossary"`) {
		t.Fatal("404 must not use bare hash nav to missing sections")
	}
	if !strings.Contains(missing, `href="./#catalog"`) || !strings.Contains(missing, `href="./#work"`) || !strings.Contains(missing, `href="./#approach"`) {
		t.Fatal("404 in-site nav should send visitors home to those sections")
	}
	if strings.Contains(missing, `href="./#glossary"`) || strings.Contains(missing, `href="#glossary"`) {
		t.Fatal("404 must not link to #glossary")
	}
	if strings.Contains(missing, `href="./" target="_blank"`) {
		t.Fatal("404 return/home must stay same-tab")
	}
	if cat.Identity.Intro != "" {
		for _, para := range splitParas(cat.Identity.Intro) {
			if !strings.Contains(home, para) {
				t.Fatal("index must render identity.intro")
			}
		}
	}
	if len(cat.Glossary) > 0 && !strings.Contains(home, cat.Glossary[0].Term) {
		t.Fatal("index must render glossary terms")
	}
	if !strings.Contains(home, `src="site.js"`) || !strings.Contains(home, `data-theme-toggle`) {
		t.Fatal("index must ship theme + section script")
	}
	first := cat.Featured[0]
	anchor := "#work-" + strings.TrimPrefix(first.ID, "product.")
	if !strings.Contains(home, `id="`+strings.TrimPrefix(anchor, "#")+`"`) {
		t.Fatal("work cards must have stable in-page ids")
	}
	if !strings.Contains(home, `href="`+anchor+`"`) {
		t.Fatal("catalog names must jump to the on-page write-up")
	}
	if strings.Contains(home, `href="`+anchor+`" target="_blank"`) {
		t.Fatal("in-page work jumps must stay same-tab")
	}
	if !strings.Contains(home, `data-accent="teal"`) {
		t.Fatal("work cards must carry product accents")
	}
}

func TestStickyNavContractInStylesheet(t *testing.T) {
	css, err := os.ReadFile(filepath.Join(repoRoot(t), "site.css"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(css)
	if !strings.Contains(text, "position: sticky") || !strings.Contains(text, "top: 0") {
		t.Fatal("site.css must pin .site-header after scroll")
	}
	if !strings.Contains(text, "scroll-padding-top") {
		t.Fatal("site.css must offset in-page jumps so headings are not under the header")
	}
	if !strings.Contains(text, "scroll-margin-top") {
		t.Fatal("site.css must give work cards scroll-margin so #work-* jumps clear the sticky header")
	}
	if !strings.Contains(text, "--header-h") {
		t.Fatal("header height token missing for scroll offset")
	}
}

func TestShippedPagesMatchRenderer(t *testing.T) {
	if err := run(repoRoot(t), true); err != nil {
		t.Fatal(err)
	}
}

type parsedAnchor struct {
	href, target, rel string
}

func assertLinkContract(t *testing.T, label, page string) {
	t.Helper()
	for _, a := range parseAnchors(page) {
		if a.href == "" {
			t.Errorf("%s: anchor missing href", label)
			continue
		}
		offsite := strings.HasPrefix(a.href, "http://") || strings.HasPrefix(a.href, "https://")
		if offsite {
			if a.target != "_blank" {
				t.Errorf("%s: %s must open in a new tab", label, a.href)
			}
			if !strings.Contains(a.rel, "noopener") {
				t.Errorf("%s: %s must include rel=noopener", label, a.href)
			}
			continue
		}
		if a.target == "_blank" {
			t.Errorf("%s: same-site %s must stay in this tab", label, a.href)
		}
	}
}

var anchorRe = regexp.MustCompile(`<a\s+([^>]+)>`)
var attrRe = regexp.MustCompile(`([A-Za-z0-9:-]+)="([^"]*)"`)

func parseAnchors(page string) []parsedAnchor {
	var out []parsedAnchor
	for _, m := range anchorRe.FindAllStringSubmatch(page, -1) {
		a := parsedAnchor{}
		for _, attr := range attrRe.FindAllStringSubmatch(m[1], -1) {
			switch attr[1] {
			case "href":
				a.href = attr[2]
			case "target":
				a.target = attr[2]
			case "rel":
				a.rel = attr[2]
			}
		}
		out = append(out, a)
	}
	return out
}
