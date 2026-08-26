package templates

import (
	"strings"
	"testing"
)

func TestCopyKnownKeys(t *testing.T) {
	if got := Copy("nav_about"); got != "About" {
		t.Fatalf("Copy(nav_about)=%q", got)
	}
	if got := Copy("home_about_page"); got != "About me" {
		t.Fatalf("Copy(home_about_page)=%q", got)
	}
	if got := Copy("home_about_kicker"); got != "Who built this" {
		t.Fatalf("Copy(home_about_kicker)=%q", got)
	}
	if got := Copy("not_found_home"); got != "Back home" {
		t.Fatalf("Copy(not_found_home)=%q", got)
	}
	if Copy("does_not_exist") != "does_not_exist" {
		t.Fatal("missing keys should echo the key")
	}
}

func TestCanonical(t *testing.T) {
	if Canonical("/") != "https://number.onkarsawarna.dev/" {
		t.Fatalf("home canonical: %s", Canonical("/"))
	}
	if Canonical("/about") != "https://number.onkarsawarna.dev/about" {
		t.Fatalf("about canonical: %s", Canonical("/about"))
	}
	if Canonical("/crossing") != "https://number.onkarsawarna.dev/crossing" {
		t.Fatalf("crossing canonical: %s", Canonical("/crossing"))
	}
}

func TestCrossingMeta(t *testing.T) {
	p := Page{Title: "The crossing", TitleKey: "title_crossing", Description: Copy("crossing_meta"), Path: "/crossing"}
	if !strings.Contains(p.MetaDescription(), "spendable pots") {
		t.Fatalf("crossing meta: %s", p.MetaDescription())
	}
	html := JSONLDCrossingHTML()
	if !strings.Contains(html, `"@type":"Article"`) || !strings.Contains(html, "/crossing") {
		t.Fatalf("crossing JSON-LD: %s", html)
	}
	if strings.Contains(html, "onkarsawarna@gmail.com") {
		t.Fatal("email must not appear in crossing JSON-LD")
	}
}

func TestJSONLDWebAppAndPerson(t *testing.T) {
	html := JSONLDHTML()
	if !strings.Contains(html, `"@type":"WebApplication"`) {
		t.Fatal("missing WebApplication")
	}
	if !strings.Contains(html, `"@type":"Person"`) {
		t.Fatal("missing Person")
	}
	if !strings.Contains(html, "Onkar Sawarna") {
		t.Fatal("missing person name")
	}
	if !strings.Contains(html, "linkedin.com/in/onkar-sawarna-569615187") {
		t.Fatal("missing LinkedIn sameAs")
	}
	if strings.Contains(html, "onkarsawarna@gmail.com") {
		t.Fatal("email must not appear in JSON-LD")
	}
}
