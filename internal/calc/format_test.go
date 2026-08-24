package calc

import "testing"

func TestFormatINR(t *testing.T) {
	if got := FormatINR(3e7); got != "₹3 Cr" {
		t.Fatalf("3e7: got %q", got)
	}
	if got := FormatINR(1_250_000); got != "₹12.5 L" {
		t.Fatalf("12.5L: got %q", got)
	}
	if got := FormatINR(5_000); got != "₹5 k" {
		t.Fatalf("5k: got %q", got)
	}
	if got := FormatINR(500); got != "₹500" {
		t.Fatalf("500: got %q", got)
	}
}

func TestFormatUSD(t *testing.T) {
	if got := FormatUSD(1_200_000); got != "$1.2M" {
		t.Fatalf("1.2M: got %q", got)
	}
	if got := FormatUSD(45_000); got != "$45k" {
		t.Fatalf("45k: got %q", got)
	}
	if got := FormatUSD(5_000); got != "$5k" {
		t.Fatalf("5k: got %q", got)
	}
	if got := FormatUSD(500); got != "$500" {
		t.Fatalf("500: got %q", got)
	}
}

func TestParseCompactMoney(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"5k", 5_000},
		{"50k", 50_000},
		{"2L", 200_000},
		{"2 lakh", 200_000},
		{"1.5Cr", 15_000_000},
		{"2M", 2_000_000},
		{"₹50,000", 50_000},
		{"1200000", 1_200_000},
	}
	for _, tc := range cases {
		got, ok := ParseCompactMoney(tc.in)
		if !ok || got != tc.want {
			t.Fatalf("%q: got %v ok=%v, want %v", tc.in, got, ok, tc.want)
		}
	}
}

func TestFormatMoney(t *testing.T) {
	if got := FormatMoney(1_250_000, "in"); got != "₹12.5 L" {
		t.Fatalf("in: got %q", got)
	}
	if got := FormatMoney(1_200_000, "us"); got != "$1.2M" {
		t.Fatalf("us: got %q", got)
	}
}
