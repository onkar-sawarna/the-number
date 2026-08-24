package calc

import (
	"math"
	"strconv"
	"strings"
)

// FormatINR formats a rupee amount: crores, lakhs, or Indian grouping.
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
	if n >= 1e4 {
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
