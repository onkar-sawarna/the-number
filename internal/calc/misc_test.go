package calc

import "testing"

func TestEmergencyTwelveMonths(t *testing.T) {
	out := Emergency(EmergencyInput{
		MonthlyEssentials: 40_000,
		MonthsCover:       6,
		CurrentBuffer:     0,
		MonthlyTopup:      20_000,
		ParkingReturn:     0,
	})
	if out.Target != 240_000 {
		t.Fatalf("target=%v", out.Target)
	}
	if !out.Reaches || out.MonthsToFill != 12 {
		t.Fatalf("months=%v reaches=%v want 12", out.MonthsToFill, out.Reaches)
	}
}

func TestBudgetTargetsAndOverspent(t *testing.T) {
	ok := Budget(BudgetInput{Income: 100_000, Needs: 50_000, Wants: 30_000, Savings: 20_000})
	if ok.TargetNeeds != 50_000 || ok.TargetWants != 30_000 || ok.TargetSavings != 20_000 {
		t.Fatalf("targets %+v", ok)
	}
	if ok.Overspent {
		t.Fatal("should not be overspent")
	}
	over := Budget(BudgetInput{Income: 100_000, Needs: 70_000, Wants: 40_000, Savings: 20_000})
	if !over.Overspent {
		t.Fatal("expected overspent flag")
	}
	if over.Unallocated >= 0 {
		t.Fatalf("unallocated=%v want negative", over.Unallocated)
	}
}

func TestHeuristicSleeves(t *testing.T) {
	for _, risk := range []string{"conservative", "moderate", "aggressive"} {
		res := HeuristicAllocation(35, 10, risk)
		sum := SleevesSum(res.Sleeves)
		if sum < 99.5 || sum > 100.5 {
			t.Fatalf("%s sum=%v", risk, sum)
		}
		for _, s := range res.Sleeves {
			for _, bad := range []string{"HDFC", "Nippon", "ticker", "AMC"} {
				if containsFold(s.Category, bad) || containsFold(res.Rationale, bad) {
					t.Fatalf("product name %q in %s", bad, risk)
				}
			}
		}
	}
	older := HeuristicAllocation(55, 2, "moderate")
	for _, s := range older.Sleeves {
		if s.Category == CatSmall && s.Percent != 0 {
			t.Fatalf("age 50+ should zero small-cap, got %v", s.Percent)
		}
	}
}

func containsFold(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && (contains(s, sub))))
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
