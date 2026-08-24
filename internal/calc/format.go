package calc

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var compactMoneyRe = regexp.MustCompile(`(?i)^(-?\d*\.?\d+)\s*(crores?|crs?|lakhs?|lacs?|l|k|m)?$`)

// ParseCompactMoney reads 50k, 2L, 1.5Cr, 2M, or a plain number.
func ParseCompactMoney(raw string) (float64, bool) {
	s := strings.TrimSpace(raw)
	s = strings.NewReplacer(",", "", "₹", "", "$", "").Replace(s)
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, true
	}
	m := compactMoneyRe.FindStringSubmatch(s)
	if m == nil {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}
	n, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	u := strings.ToLower(m[2])
	switch {
	case strings.HasPrefix(u, "cr"):
		return n * 1e7, true
	case u == "l" || strings.HasPrefix(u, "lac") || strings.HasPrefix(u, "lakh"):
		return n * 1e5, true
	case u == "k":
		return n * 1e3, true
	case u == "m":
		return n * 1e6, true
	default:
		return n, true
	}
}

// FormatINR formats a rupee amount: crores, lakhs, thousands, or Indian grouping.
func FormatINR(n float64) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	if n >= 1e7 {
		return sign + "₹" + trimZeros(n/1e7) + " Cr"
	}
	if n >= 1e5 {
		return sign + "₹" + trimZeros(n/1e5) + " L"
	}
	if n >= 1e3 {
		return sign + "₹" + trimZeros(n/1e3) + " k"
	}
	return sign + "₹" + indianGroup(int64(math.Round(n)))
}

// FormatUSD formats dollars: millions, thousands, or grouped.
func FormatUSD(n float64) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	if n >= 1e6 {
		return sign + "$" + trimZeros(n/1e6) + "M"
	}
	if n >= 1e3 {
		return sign + "$" + trimZeros(n/1e3) + "k"
	}
	return sign + "$" + westernGroup(int64(math.Round(n)))
}

func FormatMoney(n float64, region string) string {
	if region == "us" || region == "usd" {
		return FormatUSD(n)
	}
	return FormatINR(n)
}

func westernGroup(n int64) string {
	neg := ""
	if n < 0 {
		neg = "-"
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return neg + s
	}
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	if s != "" {
		parts = append([]string{s}, parts...)
	}
	return neg + strings.Join(parts, ",")
}

func trimZeros(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

func indianGroup(n int64) string {
	neg := ""
	if n < 0 {
		neg = "-"
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return neg + s
	}
	last3 := s[len(s)-3:]
	rest := s[:len(s)-3]
	var groups []string
	for len(rest) > 2 {
		groups = append([]string{rest[len(rest)-2:]}, groups...)
		rest = rest[:len(rest)-2]
	}
	if rest != "" {
		groups = append([]string{rest}, groups...)
	}
	return neg + strings.Join(groups, ",") + "," + last3
}

func monthlyEffective(annualPct float64) float64 {
	return math.Pow(1+annualPct/100.0, 1.0/12.0) - 1
}

func nearlyZero(v float64) bool {
	return math.Abs(v) < 1e-12
}
