package main

import (
	"fmt"
	"html"
	"strings"
)

// href emits an <a> opening tag. Off-site http(s) URLs open in a new tab
// with rel=noopener. Same-site targets (./, #…) stay in this tab.
func href(url string, class string) string {
	attrs := `href="` + xml(url) + `"`
	if class != "" {
		attrs += ` class="` + xml(class) + `"`
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		attrs += ` target="_blank" rel="noopener noreferrer"`
	}
	return "<a " + attrs + ">"
}

func renderIndex(cat Catalog) string {
	id := cat.Identity
	var b strings.Builder
	b.WriteString(pageHead(id, id.Name, collapseWS(id.Thesis), id.Site+"/", "#catalog", "Skip to work"))
	b.WriteString(siteHeader(id, true, len(cat.AlsoPublic) > 0))
	b.WriteString("<main id=\"main\">\n")
	b.WriteString("<section class=\"hero wrap\">\n")
	b.WriteString("<p class=\"kicker\"><span class=\"kicker-dot\" aria-hidden=\"true\"></span>Public GitHub work</p>\n")
	fmt.Fprintf(&b, "<h1>%s</h1>\n", xml(id.Name))
	fmt.Fprintf(&b, "<p class=\"lede\">%s, %s.</p>\n", xml(id.Role), xml(id.Location))
	for _, para := range splitParas(id.Intro) {
		fmt.Fprintf(&b, "<p class=\"thesis\">%s</p>\n", xml(para))
	}
	b.WriteString("<p class=\"hero-actions\">" + href("#catalog", "text-link") + "See the work</a></p>\n")
	b.WriteString("</section>\n")

	b.WriteString("<section class=\"section wrap\" id=\"catalog\">\n")
	b.WriteString("<div class=\"section-head\"><h2>Catalog</h2>")
	fmt.Fprintf(&b, "<p class=\"meta\">%d with GitHub releases</p></div>\n", len(cat.Featured))
	b.WriteString("<div class=\"catalog-board\"><table>\n")
	b.WriteString("<caption>Open a name for the write-up. Versions open GitHub.</caption>\n")
	b.WriteString("<thead><tr><th></th><th>Product</th><th>Kind</th><th>Stack</th><th>License</th><th>Release</th></tr></thead>\n<tbody>\n")
	for i, p := range cat.Featured {
		slug := strings.TrimPrefix(p.ID, "product.")
		wid := workID(p)
		detailID := "detail-" + slug
		toggleID := "toggle-" + slug
		fmt.Fprintf(&b, "<tr class=\"catalog-item\" id=\"%s\" data-accent=\"%s\">", xml(wid), xml(accentFor(p.ID)))
		fmt.Fprintf(&b, "<td data-label=\"#\" class=\"idx\">%02d</td>", i+1)
		fmt.Fprintf(&b, "<td data-label=\"Product\"><button type=\"button\" class=\"catalog-toggle\" aria-expanded=\"false\" aria-controls=\"%s\" id=\"%s\">", xml(detailID), xml(toggleID))
		fmt.Fprintf(&b, "<span class=\"product-name\">%s</span><span class=\"product-id\">%s</span></button></td>", xml(p.Name), xml(p.Repo))
		fmt.Fprintf(&b, "<td data-label=\"Kind\" class=\"lane\">%s</td>", xml(p.Lane))
		fmt.Fprintf(&b, "<td data-label=\"Stack\" class=\"stack\">%s</td>", xml(p.Language))
		fmt.Fprintf(&b, "<td data-label=\"License\"><span class=\"license\">%s</span></td>", xml(p.License))
		fmt.Fprintf(&b, "<td data-label=\"Release\">%s%s</a></td></tr>\n", href(p.ProofURL, "proof"), xml(p.Proof))

		b.WriteString("<tr class=\"catalog-detail-row\"><td colspan=\"6\">")
		fmt.Fprintf(&b, "<div class=\"catalog-detail\" id=\"%s\" role=\"region\" aria-labelledby=\"%s\" hidden>", xml(detailID), xml(toggleID))
		b.WriteString("<div class=\"catalog-detail-inner\"><div class=\"catalog-detail-body\">")
		if p.Summary != "" {
			fmt.Fprintf(&b, "<p class=\"work-summary\">%s</p>", xml(p.Summary))
		}
		fmt.Fprintf(&b, "<p>%s</p>", xml(collapseWS(p.Detail)))
		b.WriteString("<p class=\"work-meta\">")
		fmt.Fprintf(&b, "<span>%s</span><span>%s</span>", xml(p.Language), xml(p.License))
		if p.Metric.Value != "" {
			fmt.Fprintf(&b, "<span class=\"metric\">%s: %s</span>", xml(p.Metric.Label), xml(p.Metric.Value))
		}
		b.WriteString("</p>")
		if p.ID == "product.docs-puller" {
			for _, g := range cat.Glossary {
				fmt.Fprintf(&b, "<p class=\"work-term\">%s: %s</p>", xml(g.Term), xml(g.Meaning))
			}
		}
		b.WriteString("<p class=\"work-actions\">")
		fmt.Fprintf(&b, "%sSource</a>%sRelease %s</a>", href(p.URL, ""), href(p.ProofURL, ""), xml(p.Proof))
		if p.DemoURL != "" {
			label := p.DemoLabel
			if label == "" {
				label = "Live demo"
			}
			fmt.Fprintf(&b, "%s%s</a>", href(p.DemoURL, ""), xml(label))
		}
		b.WriteString("</p></div></div></div></td></tr>\n")
	}
	note := collapseWS(cat.Footnote)
	if note == "" {
		note = "Source is on GitHub."
	}
	b.WriteString("</tbody></table>\n<p class=\"catalog-note\">" + xml(note) + "</p></div>\n</section>\n")

	if len(cat.AlsoPublic) > 0 {
		b.WriteString("<section class=\"section wrap\" id=\"also-public\">\n")
		b.WriteString("<div class=\"section-head\"><h2>Also on GitHub</h2><p class=\"meta\">Public source, including repos without a tag</p></div>\n")
		b.WriteString("<ul class=\"also-public\">\n")
		for _, p := range cat.AlsoPublic {
			fmt.Fprintf(&b, "<li><strong>%s%s</a></strong> — %s <span class=\"also-meta\">%s · %s · %s%s</a></span></li>\n",
				href(p.URL, ""), xml(p.Name), xml(collapseWS(p.Summary)),
				xml(p.Language), xml(p.License), href(p.ProofURL, "proof"), xml(p.Proof))
		}
		b.WriteString("</ul>\n</section>\n")
	}

	b.WriteString("<section class=\"section wrap approach\" id=\"approach\">\n<h2>How I ship</h2>\n<div class=\"ship-stack\">\n")
	for _, p := range cat.Principles {
		fmt.Fprintf(&b, "<div class=\"ship-line\"><h3>%s</h3><p>%s</p></div>\n", xml(p.Title), xml(p.Body))
	}
	b.WriteString("</div>\n</section>\n")

	b.WriteString("<section class=\"section wrap\" id=\"toolbox\">\n<div class=\"section-head\"><h2>Toolbox</h2></div>\n<div class=\"toolbox\">")
	for _, tool := range cat.Toolbox {
		fmt.Fprintf(&b, "<span class=\"chip\">%s</span>", xml(tool))
	}
	b.WriteString("</div>\n</section>\n")
	b.WriteString("</main>\n")
	b.WriteString(siteFooter(id, cat.Footnote))
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func render404(cat Catalog) string {
	id := cat.Identity
	var b strings.Builder
	b.WriteString(pageHead(id, "Page not found", "That address is not on this site.", id.Site+"/404.html", "./", "Skip to home"))
	b.WriteString(siteHeader(id, false, false))
	b.WriteString("<main id=\"main\" class=\"missing wrap\">\n<h1>Page not found.</h1>\n")
	b.WriteString("<p>That page is not on this site. The home page lists the work.</p>\n")
	b.WriteString("<p class=\"chips\">" + href("./", "chip") + "Return home</a></p>\n</main>\n")
	b.WriteString(siteFooter(id, cat.Footnote))
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func pageHead(id Identity, title, description, canonical, skip, skipLabel string) string {
	site := strings.TrimRight(id.Site, "/")
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	b.WriteString("<meta name=\"color-scheme\" content=\"dark light\">\n")
	fmt.Fprintf(&b, "<title>%s</title>\n", xml(title))
	fmt.Fprintf(&b, "<meta name=\"description\" content=\"%s\">\n", xml(description))
	fmt.Fprintf(&b, "<link rel=\"canonical\" href=\"%s\">\n", xml(canonical))
	fmt.Fprintf(&b, "<meta property=\"og:type\" content=\"website\">\n")
	fmt.Fprintf(&b, "<meta property=\"og:title\" content=\"%s\">\n", xml(title))
	fmt.Fprintf(&b, "<meta property=\"og:description\" content=\"%s\">\n", xml(description))
	fmt.Fprintf(&b, "<meta property=\"og:url\" content=\"%s\">\n", xml(canonical))
	fmt.Fprintf(&b, "<meta property=\"og:image\" content=\"%s/assets/og.png\">\n", xml(site))
	b.WriteString("<meta name=\"twitter:card\" content=\"summary_large_image\">\n")
	b.WriteString("<meta name=\"theme-color\" content=\"#10100f\">\n")
	b.WriteString("<link rel=\"icon\" href=\"assets/mark.svg\" type=\"image/svg+xml\">\n")
	b.WriteString("<link rel=\"stylesheet\" href=\"site.css\">\n")
	b.WriteString("<noscript><style>.catalog-detail[hidden]{display:block}.catalog-detail-inner{grid-template-rows:1fr}.catalog-detail-body{opacity:1}</style></noscript>\n")
	b.WriteString("<script src=\"site.js\" defer></script>\n")
	b.WriteString("</head>\n<body>\n")
	fmt.Fprintf(&b, "<a class=\"skip\" href=\"%s\">%s</a>\n", xml(skip), xml(skipLabel))
	return b.String()
}

func siteHeader(id Identity, onHome bool, alsoPublic bool) string {
	catalog, also, approach := "#catalog", "#also-public", "#approach"
	if !onHome {
		catalog, also, approach = "./#catalog", "./#also-public", "./#approach"
	}
	var b strings.Builder
	b.WriteString("<header class=\"site-header\">\n<div class=\"wrap\">\n")
	b.WriteString(href("./", "brand") + markSVG() + "<span class=\"brand-name\">" + xml(id.Name) + "</span></a>\n")
	b.WriteString("<nav aria-label=\"Primary\">\n")
	fmt.Fprintf(&b, "%s</a>\n", navItem(catalog, "catalog", "Catalog"))
	if alsoPublic {
		fmt.Fprintf(&b, "%s</a>\n", navItem(also, "also-public", "Also on GitHub"))
	}
	fmt.Fprintf(&b, "%s</a>\n%sGitHub</a>\n",
		navItem(approach, "approach", "Approach"),
		href(id.GitHub, ""))
	b.WriteString("</nav>\n")
	b.WriteString("<button type=\"button\" class=\"theme-toggle\" data-theme-toggle aria-label=\"Color theme: System. Switch theme\">System</button>\n")
	b.WriteString("</div>\n</header>\n")
	return b.String()
}

func siteFooter(id Identity, note string) string {
	if collapseWS(note) == "" {
		note = "Source is on GitHub."
	}
	return "<footer class=\"site-footer\">\n<div class=\"wrap\">\n<div><p class=\"brand-name\">" + xml(id.Name) + "</p><p class=\"boundary\">" + xml(collapseWS(note)) + "</p></div>\n<div class=\"footer-links\">" + href(id.GitHub, "") + "GitHub</a>" + href(id.LinkedIn, "") + "LinkedIn</a></div>\n</div>\n</footer>\n"
}

func markSVG() string {
	return `<svg viewBox="0 0 79 79" role="img" aria-label="Catalog mark"><title>Catalog N</title><rect x="0" y="0" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="21" y="0" width="16" height="16" rx="3.5" fill="currentColor" opacity=".12"/><rect x="42" y="0" width="16" height="16" rx="3.5" fill="currentColor" opacity=".12"/><rect x="63" y="0" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="0" y="21" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="21" y="21" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="42" y="21" width="16" height="16" rx="3.5" fill="currentColor" opacity=".12"/><rect x="63" y="21" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="0" y="42" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="21" y="42" width="16" height="16" rx="3.5" fill="currentColor" opacity=".12"/><rect x="42" y="42" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="63" y="42" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="0" y="63" width="16" height="16" rx="3.5" fill="#5c7cff"/><rect x="21" y="63" width="16" height="16" rx="3.5" fill="currentColor" opacity=".12"/><rect x="42" y="63" width="16" height="16" rx="3.5" fill="currentColor" opacity=".12"/><rect x="63" y="63" width="16" height="16" rx="3.5" fill="#5c7cff"/></svg>`
}

func renderMark() string {
	cells := ""
	on := map[[2]int]bool{
		{0, 0}: true, {3, 0}: true,
		{0, 1}: true, {1, 1}: true, {3, 1}: true,
		{0, 2}: true, {2, 2}: true, {3, 2}: true,
		{0, 3}: true, {3, 3}: true,
	}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			x, y := col*21, row*21
			if on[[2]int{col, row}] {
				cells += fmt.Sprintf(`<rect x="%d" y="%d" width="16" height="16" rx="3.5" fill="#5c7cff"/>`, x, y)
			} else {
				cells += fmt.Sprintf(`<rect x="%d" y="%d" width="16" height="16" rx="3.5" fill="#10100f"/>`, x, y)
			}
		}
	}
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 79 79" role="img"><title>Catalog N</title>` + cells + `</svg>` + "\n"
}

func renderRobots(cat Catalog) string {
	return "User-agent: *\nAllow: /\nSitemap: " + strings.TrimRight(cat.Identity.Site, "/") + "/\n"
}

func navItem(url, section, label string) string {
	return fmt.Sprintf(`<a href="%s" data-section="%s">%s`, xml(url), xml(section), xml(label))
}

func workID(p Product) string {
	return "work-" + strings.TrimPrefix(p.ID, "product.")
}

func accentFor(id string) string {
	switch id {
	case "product.docs-puller":
		return "teal"
	case "product.nicos-catalog":
		return "cobalt"
	case "product.openbook":
		return "ember"
	case "product.agent-ops":
		return "mint"
	case "product.nicos-hidden-menubar":
		return "violet"
	case "product.nicos-slot-dock":
		return "steel"
	case "product.jobkit":
		return "gold"
	case "product.session-pressure":
		return "mint"
	case "product.keepawake":
		return "ember"
	default:
		return "cobalt"
	}
}

func xml(s string) string {
	return html.EscapeString(s)
}

func collapseWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func splitParas(s string) []string {
	var out []string
	for _, para := range strings.Split(s, "\n\n") {
		para = collapseWS(para)
		if para != "" {
			out = append(out, para)
		}
	}
	if len(out) == 0 && collapseWS(s) != "" {
		return []string{collapseWS(s)}
	}
	return out
}
