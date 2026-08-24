package calc

import "testing"

func TestFormatINR(t *testing.T) {
	if got := FormatINR(3e7); got != "₹3 Cr" {
		t.Fatalf("3e7: got %q", got)
	}
	if got := FormatINR(1_250_000); got != "₹12.5 L" {
		t.Fatalf("12.5L: got %q", got)
	}
	if got := FormatINR(12_345); got != "₹12,345" {
		t.Fatalf("12345: got %q", got)
	}
}

func TestFormatUSD(t *testing.T) {
	if got := FormatUSD(1_200_000); got != "$1.2M" {
		t.Fatalf("1.2M: got %q", got)
	}
	if got := FormatUSD(45_000); got != "$45k" {
		t.Fatalf("45k: got %q", got)
	}
	if got := FormatUSD(1234); got != "$1,234" {
		t.Fatalf("1234: got %q", got)
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
