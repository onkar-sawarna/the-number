package templates

import (
	"net/url"
	"os"
	"strings"

	"github.com/a-h/templ"
)

// DefaultSupportUPI is shown on Support and after FIRE results when SUPPORT_UPI is unset.
// Override on the host if this VPA is not the one you use.
const DefaultSupportUPI = "onkarsawarna-3@okicici"

// SupportUPINote is the UPI transaction remark (tn) so the credit shows as FIRE.
const SupportUPINote = "Paid for FIRE — the number"

// SupportUPIQRSrc is the cropped Google Pay QR (code only) for DefaultSupportUPI.
const SupportUPIQRSrc = "/static/img/upi-qr-code.png?v=1"

func SupportUPI() string {
	if v := strings.TrimSpace(os.Getenv("SUPPORT_UPI")); v != "" {
		return v
	}
	return DefaultSupportUPI
}

func SupportUPIPayString() string {
	q := url.Values{}
	q.Set("pa", SupportUPI())
	q.Set("pn", "Onkar Sawarna")
	q.Set("cu", "INR")
	q.Set("tn", SupportUPINote)
	return "upi://pay?" + q.Encode()
}

func SupportUPIHref() templ.SafeURL {
	return templ.SafeURL(SupportUPIPayString())
}

// DefaultSupportCardURL is the World (non-UPI) coffee page when SUPPORT_CARD_URL is unset.
// Claim this exact Buy Me a Coffee username or set SUPPORT_CARD_URL to Ko-fi / PayPal.me.
const DefaultSupportCardURL = "https://buymeacoffee.com/onkarsawarna"

func SupportCardURL() string {
	if v := strings.TrimSpace(os.Getenv("SUPPORT_CARD_URL")); v != "" {
		return v
	}
	return DefaultSupportCardURL
}

func supportCardAllowed(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" {
		return false
	}
	switch strings.ToLower(u.Hostname()) {
	case "buymeacoffee.com", "www.buymeacoffee.com", "ko-fi.com", "www.ko-fi.com", "paypal.me", "www.paypal.me":
		return true
	case "github.com":
		return strings.HasPrefix(u.Path, "/sponsors/")
	default:
		return false
	}
}

func SupportCardHref() templ.SafeURL {
	u := SupportCardURL()
	if !supportCardAllowed(u) {
		u = DefaultSupportCardURL
	}
	return templ.SafeURL(u)
}
