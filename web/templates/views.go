package templates

import (
	"strconv"
	"strings"
)

type Page struct {
	Title    string
	Dark     bool
	LoggedIn bool
	Email    string
	Path     string
}

type AuthView struct {
	Page
	Error string
	Value string
	Next  string
}

type ScenarioVM struct {
	ID         uint
	Date       string
	Kind       string
	KindLabel  string
	Title      string
	DeletePath string
}

type KV struct {
	Label string
	Value string
}

type CompareCard struct {
	Title string
	Kind  string
	Date  string
	Rows  []KV
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

func uintStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func TabClass(path, href string) string {
	base := "rounded-full px-3 py-1.5 text-sm transition-colors whitespace-nowrap "
	if path == href {
		return base + "bg-brand text-white"
	}
	return base + "text-ink/70 hover:bg-black/5 dark:text-zinc-300 dark:hover:bg-white/10"
}
