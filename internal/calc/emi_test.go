package calc

import (
	"math"
	"testing"
)

func TestEMITenLakh(t *testing.T) {
	emi := EMIAmount(1_000_000, 10, 20)
	if math.Abs(emi-9650) > 20 {
		t.Fatalf("EMI=%v want ≈9650", emi)
	}
}

func TestEMIPrepaymentSaves(t *testing.T) {
	base := EMI(EMIInput{Principal: 1_000_000, AnnualRate: 10, Years: 20, ExtraMonthly: 0, Lump: 0})
	extra := EMI(EMIInput{Principal: 1_000_000, AnnualRate: 10, Years: 20, ExtraMonthly: 2_000, Lump: 50_000})
	if extra.InterestSaved <= 0 {
		t.Fatalf("expected interest saved, got %v", extra.InterestSaved)
	}
	if extra.MonthsSaved <= 0 {
		t.Fatalf("expected months saved, got %v (base months %v prepaid %v)", extra.MonthsSaved, base.Months, extra.PrepaidMonths)
	}
	if extra.PrepaidInterest >= base.TotalInterest {
		t.Fatal("prepaid interest should be lower")
	}
}
