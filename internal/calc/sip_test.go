package calc

import "testing"

func TestSIPTenYears(t *testing.T) {
	out := SIP(SIPInput{Monthly: 10_000, Existing: 0, ExpectedReturn: 12, Years: 10})
	if out.Invested != 1_200_000 {
		t.Fatalf("invested=%v want 1200000", out.Invested)
	}
	if out.FV < 2_200_000 || out.FV > 2_500_000 {
		t.Fatalf("FV=%v want roughly 22–25L", out.FV)
	}
}

func TestSIPZeroRate(t *testing.T) {
	out := SIP(SIPInput{Monthly: 10_000, Existing: 50_000, ExpectedReturn: 0, Years: 10})
	want := 50_000.0 + 10_000*120
	if out.FV != want {
		t.Fatalf("zero rate FV=%v want %v", out.FV, want)
	}
}
