package calc

import "testing"

func TestFIRENumberLeanFat(t *testing.T) {
	expenses := 1_200_000.0
	swr := 4.0
	n := FIRENumber(expenses, swr)
	if n != 30_000_000 {
		t.Fatalf("FIRE number: got %v want 30000000", n)
	}
	out := FIRE(FIREInput{AnnualExpenses: expenses, SWR: swr, Age: 30})
	if out.Lean != n*0.5 || out.Fat != n*2 || out.Regular != n {
		t.Fatalf("lean/regular/fat: %+v", out)
	}
}

func TestAlreadyFIZeroYears(t *testing.T) {
	in := FIREInput{
		Age:            40,
		AnnualExpenses: 1_200_000,
		CurrentSavings: 40_000_000,
		MonthlySavings: 0,
		ExpectedReturn: 8,
		Inflation:      6,
		SWR:            4,
	}
	out := FIRE(in)
	if !out.ReachesFire || out.Years != 0 {
		t.Fatalf("already FI: years=%v reaches=%v", out.Years, out.ReachesFire)
	}
}

func TestZeroReturnInflationTenYears(t *testing.T) {
	// 25k/mo vs 30L target, 0% return, 0% inflation → ~10 years
	in := FIREInput{
		Age:            30,
		AnnualExpenses: 120_000, // 30L target at 4% SWR
		CurrentSavings: 0,
		MonthlySavings: 25_000,
		ExpectedReturn: 0,
		Inflation:      0,
		SWR:            4,
	}
	out := FIRE(in)
	if !out.ReachesFire {
		t.Fatal("expected to reach FIRE")
	}
	if out.Years < 9.5 || out.Years > 10.5 {
		t.Fatalf("years=%v want ~10", out.Years)
	}
}

func TestHigherMonthlyFewerYears(t *testing.T) {
	base := DefaultFIRE()
	low := FIRE(base)
	base.MonthlySavings = base.MonthlySavings * 2
	high := FIRE(base)
	if !low.ReachesFire || !high.ReachesFire {
		t.Fatal("both should reach")
	}
	if high.Years >= low.Years {
		t.Fatalf("higher savings should be faster: %v vs %v", high.Years, low.Years)
	}
}

func TestZeroSavingsZeroReturnNever(t *testing.T) {
	in := FIREInput{
		Age:            30,
		AnnualExpenses: 1_200_000,
		CurrentSavings: 0,
		MonthlySavings: 0,
		ExpectedReturn: 0,
		Inflation:      0,
		SWR:            4,
	}
	out := FIRE(in)
	if out.ReachesFire {
		t.Fatalf("should never reach, got years=%v", out.Years)
	}
}
