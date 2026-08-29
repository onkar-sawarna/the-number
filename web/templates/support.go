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
