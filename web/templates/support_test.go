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
	if !strings.Contains(SupportUPIQRSrc, "/static/img/upi-qr-code.png") {
		t.Fatalf("SupportUPIQRSrc=%q", SupportUPIQRSrc)
	}
}

func TestSupportCardDefault(t *testing.T) {
	t.Setenv("SUPPORT_CARD_URL", "")
	if SupportCardURL() != DefaultSupportCardURL {
		t.Fatalf("SupportCardURL()=%q", SupportCardURL())
	}
	if string(SupportCardHref()) != DefaultSupportCardURL {
		t.Fatalf("SupportCardHref()=%q", SupportCardHref())
	}
}

func TestSupportCardEnv(t *testing.T) {
	t.Setenv("SUPPORT_CARD_URL", "https://ko-fi.com/onkarsawarna")
	if SupportCardURL() != "https://ko-fi.com/onkarsawarna" {
		t.Fatalf("SupportCardURL()=%q", SupportCardURL())
	}
	if string(SupportCardHref()) != "https://ko-fi.com/onkarsawarna" {
		t.Fatalf("SupportCardHref()=%q", SupportCardHref())
	}
}

func TestSupportCardRejectsJunk(t *testing.T) {
	t.Setenv("SUPPORT_CARD_URL", "javascript:alert(1)")
	if string(SupportCardHref()) != DefaultSupportCardURL {
		t.Fatalf("unsafe href leaked: %s", SupportCardHref())
	}
}
