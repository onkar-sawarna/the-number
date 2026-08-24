package ai

import (
	"strings"
	"testing"

	"github.com/thenumber/app/internal/calc"
)

func TestParseAcceptsCategorySleeves(t *testing.T) {
	raw := `{"sleeves":[{"category":"Large-cap index","percent":60},{"category":"Debt / government securities","percent":40}],"rationale":"A simple split."}`
	out, err := ParseGuidanceJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Sleeves) != 2 {
		t.Fatalf("sleeves=%d", len(out.Sleeves))
	}
	sum := calc.SleevesSum(out.Sleeves)
	if sum < 99.5 || sum > 100.5 {
		t.Fatalf("sum=%v", sum)
	}
}

func TestParseDropsProductNames(t *testing.T) {
	raw := `{"sleeves":[{"category":"HDFC Large-cap index","percent":60},{"category":"Nippon Gold","percent":10},{"category":"Debt / government securities","percent":30}],"rationale":"HDFC is not advice; stay with categories."}`
	out, err := ParseGuidanceJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if HasProductName(out.Rationale) {
		t.Fatalf("rationale still has product: %q", out.Rationale)
	}
	for _, s := range out.Sleeves {
		if HasProductName(s.Category) {
			t.Fatalf("category still has product: %q", s.Category)
		}
		if strings.Contains(strings.ToLower(s.Category), "hdfc") {
			t.Fatalf("HDFC leaked: %q", s.Category)
		}
	}
	foundLarge := false
	for _, s := range out.Sleeves {
		if s.Category == calc.CatLarge {
			foundLarge = true
		}
	}
	if !foundLarge {
		t.Fatalf("expected mapped large-cap, got %+v", out.Sleeves)
	}
}

func TestParseRejectsOnlyProducts(t *testing.T) {
	raw := `{"sleeves":[{"category":"HDFC Flexi Cap","percent":100}],"rationale":"Buy HDFC"}`
	out, err := ParseGuidanceJSON(raw)
	if err == nil && len(out.Sleeves) > 0 {
		for _, s := range out.Sleeves {
			if HasProductName(s.Category) {
				t.Fatalf("kept product sleeve %q", s.Category)
			}
		}
	}
}
