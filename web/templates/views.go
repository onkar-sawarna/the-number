package templates

import "strings"

type Page struct {
	Title string
	Dark  bool
	Path  string
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
	base := "rounded-full px-3 py-1.5 text-sm transition-colors "
	if active {
		return base + "bg-ink text-paper dark:bg-zinc-100 dark:text-night"
	}
	return base + "text-ink/70 hover:text-ink dark:text-zinc-300 dark:hover:text-white"
}

func TabClass(path, href string) string {
	base := "rounded-full px-3 py-1.5 text-sm transition-colors whitespace-nowrap "
	if path == href {
		return base + "bg-brand text-white"
	}
	return base + "text-ink/70 hover:bg-black/5 dark:text-zinc-300 dark:hover:bg-white/10"
}
