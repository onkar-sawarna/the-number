package templates

import "testing"

func TestCopyKnownKeys(t *testing.T) {
	if got := Copy("nav_about"); got != "About" {
		t.Fatalf("Copy(nav_about)=%q", got)
	}
	if got := Copy("home_about_page"); got != "About me" {
		t.Fatalf("Copy(home_about_page)=%q", got)
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
