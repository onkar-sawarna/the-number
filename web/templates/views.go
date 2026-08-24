package templates

import "strings"

type Page struct {
	Title    string
	TitleKey string
	Dark     bool
	Path     string
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
