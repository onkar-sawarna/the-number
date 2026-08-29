package templates

import (
	"strings"
	"testing"
)

func TestSupportUPIDefault(t *testing.T) {
	t.Setenv("SUPPORT_UPI", "")
	if SupportUPI() != DefaultSupportUPI {
		t.Fatalf("SupportUPI()=%q", SupportUPI())
	}
	href := string(SupportUPIHref())
	if !strings.HasPrefix(href, "upi://pay?") || !strings.Contains(href, "pa=") {
		t.Fatalf("SupportUPIHref()=%q", href)
	}
}

func TestSupportUPIEnv(t *testing.T) {
	t.Setenv("SUPPORT_UPI", "onkar@okaxis")
	if SupportUPI() != "onkar@okaxis" {
		t.Fatalf("SupportUPI()=%q", SupportUPI())
	}
	if !strings.Contains(string(SupportUPIHref()), "onkar") {
		t.Fatalf("href missing env VPA: %s", SupportUPIHref())
	}
}
