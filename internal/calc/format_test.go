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
