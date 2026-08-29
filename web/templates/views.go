package templates

import (
	"encoding/json"
	"strings"
)

const SiteOrigin = "https://number.onkarsawarna.dev"

type Page struct {
	Title       string
	TitleKey    string
	Description string
	Dark        bool
	Path        string
}

func (p Page) MetaDescription() string {
	if strings.TrimSpace(p.Description) != "" {
		return p.Description
	}
	return "FIRE is Financial Independence, Retire Early — a corpus that can cover your expenses without a salary. Know the number, see the year you could stop. No login, no fund pitch."
}

type SleeveVM struct {
	Category string
	Percent  string
}

type GuidanceResultView struct {
	Warning   string
	Rationale string
	Source    string
	Labels    string
	Values    string
	Sleeves   []SleeveVM
}

func MainNavClass(path, href string) string {
	active := path == href
	if href == "/calculators/fire" && strings.HasPrefix(path, "/calculators") {
		active = true
	}
	if href == "/" {
		active = path == "/"
	}
	base := "inline-flex min-h-10 shrink-0 items-center rounded-full px-2.5 py-1.5 text-sm transition-colors sm:px-3 "
	if active {
		return base + "bg-ink text-paper dark:bg-zinc-100 dark:text-night"
	}
	return base + "text-ink/70 hover:bg-black/5 hover:text-ink dark:text-zinc-300 dark:hover:bg-white/10 dark:hover:text-white"
}

func DockClass(path, href string) string {
	active := path == href
	if href == "/calculators/fire" && strings.HasPrefix(path, "/calculators") {
		active = true
	}
	if href == "/" {
		active = path == "/"
	}
	base := "flex min-h-12 flex-col items-center justify-center gap-0.5 text-[11px] font-medium "
	if active {
		return base + "text-brand"
	}
	return base + "text-ink/50 dark:text-zinc-400"
}

func IsCalc(path string) bool {
	return strings.HasPrefix(path, "/calculators")
}

func TabClass(path, href string) string {
	base := "inline-flex min-h-10 shrink-0 items-center rounded-full px-3 py-1.5 text-sm transition-colors whitespace-nowrap "
	if path == href {
		return base + "bg-brand text-white"
	}
	return base + "text-ink/70 hover:bg-black/5 dark:text-zinc-300 dark:hover:bg-white/10"
}

func Canonical(path string) string {
	if path == "" || path == "/" {
		return SiteOrigin + "/"
	}
	return SiteOrigin + path
}

func OgImageURL() string {
	return SiteOrigin + "/static/img/og.png"
}

func JSONLDHTML() string {
	return `<script type="application/ld+json">` + jsonLDDocument() + `</script>`
}

func jsonLDDocument() string {
	doc := map[string]any{
		"@context": "https://schema.org",
		"@graph": []any{
			map[string]any{
				"@type":               "WebApplication",
				"@id":                 SiteOrigin + "/#app",
				"name":                "the number",
				"url":                 SiteOrigin + "/",
				"description":         "A FIRE calculator that explains the number, then lets you play with it. SIP, EMI, and budget tools too. Numbers stay in your browser.",
				"applicationCategory": "FinanceApplication",
				"operatingSystem":     "Any",
				"browserRequirements": "Requires JavaScript",
				"image":               OgImageURL(),
				"offers":              map[string]any{"@type": "Offer", "price": "0", "priceCurrency": "INR"},
				"author":              map[string]any{"@id": SiteOrigin + "/#person"},
			},
			map[string]any{
				"@type": "Person",
				"@id":   SiteOrigin + "/#person",
				"name":  "Onkar Sawarna",
				"url":   SiteOrigin + "/about",
				"sameAs": []string{
					"https://github.com/onkar-sawarna",
					"https://www.linkedin.com/in/onkar-sawarna-569615187/",
					"https://x.com/onkar_sawarna",
					"https://www.onkarsawarna.dev/",
				},
			},
		},
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func JSONLDCrossingHTML() string {
	doc := map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Article",
		"@id":         SiteOrigin + "/crossing#article",
		"headline":    "The crossing",
		"name":        "The crossing",
		"description": "The crossing is the first year your spendable pots earn more than you add. It is not FIRE. Jewellery you keep is out.",
		"url":         SiteOrigin + "/crossing",
		"isPartOf":    map[string]any{"@id": SiteOrigin + "/#app"},
		"author":      map[string]any{"@id": SiteOrigin + "/#person"},
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return ""
	}
	return `<script type="application/ld+json">` + string(b) + `</script>`
}
