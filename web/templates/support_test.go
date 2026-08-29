package templates

import (
	"net/url"
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
	decoded, err := url.QueryUnescape(strings.ReplaceAll(href, "+", " "))
	if err != nil || !strings.Contains(href, "tn=") || !strings.Contains(decoded, SupportUPINote) {
		t.Fatalf("pay link missing FIRE note: %s", href)
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

func TestSupportUPIQR(t *testing.T) {
	if !strings.Contains(SupportUPIQRSrc, "/static/img/upi-qr.png") {
		t.Fatalf("SupportUPIQRSrc=%q", SupportUPIQRSrc)
	}
}
