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
