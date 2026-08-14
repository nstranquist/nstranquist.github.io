package main

import (
	"fmt"
	"html"
	"strings"
)

func renderIndex(cat Catalog) string {
	id := cat.Identity
	var b strings.Builder
	b.WriteString(pageHead(id, id.Name+" — public catalog", collapseWS(id.Thesis), id.Site+"/"))
	b.WriteString(siteHeader(id, false))
	b.WriteString("<main>\n")
	b.WriteString(`<section class="hero wrap">`)
	b.WriteString(`<p class="kicker"><span class="kicker-dot" aria-hidden="true"></span>Public catalog · inspectable source</p>`)
	fmt.Fprintf(&b, "<h1>%s</h1>\n", xml(id.Name))
	fmt.Fprintf(&b, `<p class="lede">%s, %s.</p>`+"\n", xml(id.Role), xml(id.Location))
	fmt.Fprintf(&b, `<p class="thesis">I build %s. This site is the inspectable extract of that practice. Each product has a license, a proof tag, and a claim boundary.</p>`+"\n", xml(uncapitalize(strings.TrimSuffix(id.Thesis, "."))))
	b.WriteString(`<div class="chips"><span class="chip">Platforms</span><span class="chip">AI infrastructure</span><span class="chip">Full-stack</span></div>`)
	b.WriteString("</section>\n")

	b.WriteString(`<section class="section wrap" id="catalog">`)
	b.WriteString(`<div class="section-head"><h2>Catalog</h2>`)
	fmt.Fprintf(&b, `<p class="meta">%d products · ready</p></div>`+"\n", len(cat.Featured))
	b.WriteString(`<div class="catalog-board"><table>`)
	b.WriteString(`<caption>Public catalog of inspectable products</caption>`)
	b.WriteString(`<thead><tr><th>Product</th><th>Lane</th><th>Stack</th><th>License</th><th>Proof</th></tr></thead><tbody>`)
	for _, p := range cat.Featured {
		fmt.Fprintf(&b, `<tr><td data-label="Product"><a href="%s"><span class="product-name">%s</span><span class="product-id">%s</span></a></td>`, xml(p.URL), xml(p.Name), xml(p.ID))
		fmt.Fprintf(&b, `<td data-label="Lane" class="lane">%s</td>`, xml(p.Lane))
		fmt.Fprintf(&b, `<td data-label="Stack" class="stack">%s</td>`, xml(p.Language))
		fmt.Fprintf(&b, `<td data-label="License"><span class="license">%s</span></td>`, xml(p.License))
		fmt.Fprintf(&b, `<td data-label="Proof"><a class="proof" href="%s">%s</a></td></tr>`+"\n", xml(p.ProofURL), xml(p.Proof))
	}
	b.WriteString(`</tbody></table><p class="catalog-note">boundary — public extracts only · claims fail closed · no invented metrics</p></div></section>`)

	b.WriteString(`<section class="section wrap work" id="work"><h2>Selected work</h2><div class="work-list">`)
	for _, p := range cat.Featured {
		b.WriteString(`<article class="work-card">`)
		fmt.Fprintf(&b, `<header><h3><a href="%s">%s</a></h3><span class="meta">%s</span></header>`, xml(p.URL), xml(p.Name), xml(p.Lane))
		fmt.Fprintf(&b, `<p>%s</p>`, xml(collapseWS(p.Detail)))
		b.WriteString(`<p class="work-meta">`)
		fmt.Fprintf(&b, `<span>%s</span><span>%s</span><a href="%s">%s</a>`, xml(p.Language), xml(p.License), xml(p.ProofURL), xml(p.Proof))
		if p.Metric.Value != "" {
			fmt.Fprintf(&b, `<span class="metric">%s: %s</span>`, xml(p.Metric.Label), xml(p.Metric.Value))
		}
		b.WriteString(`</p></article>`)
	}
	b.WriteString(`</div></section>`)

	b.WriteString(`<section class="section wrap approach" id="approach"><h2>How I work</h2><div class="approach-grid">`)
	for _, p := range cat.Principles {
		fmt.Fprintf(&b, `<article class="principle"><h3>%s</h3><p>%s</p></article>`, xml(p.Title), xml(p.Body))
	}
	b.WriteString(`</div></section>`)

	b.WriteString(`<section class="section wrap" id="toolbox"><div class="section-head"><h2>Toolbox</h2></div><div class="toolbox">`)
	for _, tool := range cat.Toolbox {
		fmt.Fprintf(&b, `<span class="chip">%s</span>`, xml(tool))
	}
	b.WriteString(`</div></section></main>`)
	b.WriteString(siteFooter(id))
	b.WriteString("</body></html>\n")
	return b.String()
}

func render404(cat Catalog) string {
	id := cat.Identity
	var b strings.Builder
	b.WriteString(pageHead(id, "Not in this catalog", "This path is not part of the public catalog.", id.Site+"/404.html"))
	b.WriteString(siteHeader(id, false))
	b.WriteString(`<main class="missing wrap"><h1>Not in this catalog.</h1>`)
	b.WriteString(`<p>That path is not part of the public extract. The six featured products live on the home page.</p>`)
	b.WriteString(`<p class="chips"><a class="chip" href="./">Return home</a></p></main>`)
	b.WriteString(siteFooter(id))
	b.WriteString("</body></html>\n")
	return b.String()
}

func pageHead(id Identity, title, description, canonical string) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	fmt.Fprintf(&b, "<title>%s</title>\n", xml(title))
	fmt.Fprintf(&b, "<meta name=\"description\" content=\"%s\">\n", xml(description))
	fmt.Fprintf(&b, "<link rel=\"canonical\" href=\"%s\">\n", xml(canonical))
	fmt.Fprintf(&b, "<meta property=\"og:title\" content=\"%s\">\n", xml(title))
	fmt.Fprintf(&b, "<meta property=\"og:description\" content=\"%s\">\n", xml(description))
	fmt.Fprintf(&b, "<meta property=\"og:url\" content=\"%s\">\n", xml(canonical))
	fmt.Fprintf(&b, "<meta property=\"og:image\" content=\"%s/assets/og.png\">\n", xml(strings.TrimRight(id.Site, "/")))
	b.WriteString("<meta name=\"theme-color\" content=\"#10100f\">\n")
	b.WriteString("<link rel=\"icon\" href=\"assets/mark.svg\" type=\"image/svg+xml\">\n")
	b.WriteString("<link rel=\"stylesheet\" href=\"site.css\">\n")
	b.WriteString("</head>\n<body>\n")
	b.WriteString(`<a class="skip" href="#work">Skip to work</a>` + "\n")
	return b.String()
}

func siteHeader(id Identity, _ bool) string {
	return fmt.Sprintf(`<header class="site-header"><div class="wrap"><a class="brand" href="./">%s<span class="brand-name">Public catalog</span></a><nav><a href="#catalog">Catalog</a><a href="#work">Work</a><a href="#approach">Approach</a><a href="%s">GitHub</a></nav></div></header>`+"\n", markSVG(), xml(id.GitHub))
}

func siteFooter(id Identity) string {
	return fmt.Sprintf(`<footer class="site-footer"><div class="wrap"><div><p class="brand-name">%s</p><p class="boundary">public extracts only · claims fail closed · no invented metrics</p></div><div class="footer-links"><a href="%s">GitHub</a><a href="%s">LinkedIn</a></div></div></footer>`+"\n", xml(id.Name), xml(id.GitHub), xml(id.LinkedIn))
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

func xml(s string) string {
	return html.EscapeString(s)
}

func collapseWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func uncapitalize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
