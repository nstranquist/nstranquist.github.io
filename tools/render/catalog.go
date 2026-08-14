package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Catalog struct {
	SchemaVersion int             `yaml:"schema_version"`
	Identity      Identity        `yaml:"identity"`
	Featured      []Product       `yaml:"featured"`
	Toolbox       []string        `yaml:"toolbox"`
	Principles    []Principle     `yaml:"principles"`
	Footnote      string          `yaml:"footnote"`
	Glossary      []GlossaryEntry `yaml:"glossary"`
}

type Identity struct {
	Name     string `yaml:"name"`
	Role     string `yaml:"role"`
	Location string `yaml:"location"`
	Thesis   string `yaml:"thesis"`
	Intro    string `yaml:"intro"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Site     string `yaml:"site"`
}

type GlossaryEntry struct {
	Term    string `yaml:"term"`
	Meaning string `yaml:"meaning"`
}

type Product struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	Repo     string `yaml:"repo"`
	URL      string `yaml:"url"`
	ProofURL string `yaml:"proof_url"`
	License  string `yaml:"license"`
	Language string `yaml:"language"`
	Lane     string `yaml:"lane"`
	Proof    string `yaml:"proof"`
	Summary  string `yaml:"summary"`
	Detail   string `yaml:"detail"`
	Metric   Metric `yaml:"metric"`
}

type Metric struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

type Principle struct {
	Title string `yaml:"title"`
	Body  string `yaml:"body"`
}

func loadCatalog(path string) (Catalog, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Catalog{}, err
	}
	var cat Catalog
	if err := yaml.Unmarshal(raw, &cat); err != nil {
		return Catalog{}, err
	}
	if err := cat.validate(); err != nil {
		return Catalog{}, err
	}
	return cat, nil
}

func (c Catalog) validate() error {
	if c.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema_version %d", c.SchemaVersion)
	}
	if strings.TrimSpace(c.Identity.Name) == "" || strings.TrimSpace(c.Identity.Site) == "" {
		return fmt.Errorf("identity.name and identity.site are required")
	}
	if strings.TrimSpace(c.Identity.Intro) == "" {
		return fmt.Errorf("identity.intro is required")
	}
	if len(c.Featured) == 0 {
		return fmt.Errorf("featured catalog is empty")
	}
	if len(c.Glossary) == 0 {
		return fmt.Errorf("glossary is required")
	}
	seen := map[string]struct{}{}
	for i, p := range c.Featured {
		if p.ID == "" || p.Name == "" || p.URL == "" || p.ProofURL == "" {
			return fmt.Errorf("featured[%d] is missing required fields", i)
		}
		if !strings.HasPrefix(p.ID, "product.") {
			return fmt.Errorf("featured[%d].id %q must start with product.", i, p.ID)
		}
		if !strings.HasPrefix(p.URL, "https://github.com/nstranquist/") {
			return fmt.Errorf("featured[%d].url is not a public nstranquist GitHub URL", i)
		}
		if !strings.HasPrefix(p.ProofURL, p.URL+"/releases/tag/") {
			return fmt.Errorf("featured[%d].proof_url must be a release tag on %s", i, p.URL)
		}
		if _, ok := seen[p.ID]; ok {
			return fmt.Errorf("duplicate featured id %s", p.ID)
		}
		seen[p.ID] = struct{}{}
	}
	return nil
}
