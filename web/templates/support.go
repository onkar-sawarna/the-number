package templates

import (
	"net/url"
	"os"
	"strings"

	"github.com/a-h/templ"
)

// DefaultSupportUPI is shown after FIRE results when SUPPORT_UPI is unset.
// Override on the host if this VPA is not the one you use.
const DefaultSupportUPI = "onkarsawarna@oksbi"

func SupportUPI() string {
	if v := strings.TrimSpace(os.Getenv("SUPPORT_UPI")); v != "" {
		return v
	}
	return DefaultSupportUPI
}

func SupportUPIHref() templ.SafeURL {
	q := url.Values{}
	q.Set("pa", SupportUPI())
	q.Set("pn", "Onkar Sawarna")
	q.Set("cu", "INR")
	q.Set("tn", "the number")
	return templ.SafeURL("upi://pay?" + q.Encode())
}
