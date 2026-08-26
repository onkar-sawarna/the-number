package calc

import "testing"

func TestSipBumpShortensYears(t *testing.T) {
	in := DefaultFIRE()
	base := FIRE(in)
	moves := YearMoves(in)
	if len(moves) != 3 {
		t.Fatalf("moves=%d", len(moves))
	}
	sip := moves[0]
	if sip.Kind != "sip" || sip.Amount != 5_000 {
		t.Fatalf("sip move: %+v", sip)
	}
	if !sip.Reaches || sip.Years >= base.Years {
		t.Fatalf("extra SIP should be faster: base=%v sip=%v", base.Years, sip.Years)
	}
	if sip.DeltaYears >= 0 {
		t.Fatalf("delta should be negative (earlier): %v", sip.DeltaYears)
	}
}

func TestLowerReturnLengthensYears(t *testing.T) {
	in := DefaultFIRE()
	base := FIRE(in)
	ret := YearMoves(in)[1]
	if ret.Kind != "return" {
		t.Fatalf("kind=%s", ret.Kind)
	}
	if !ret.Reaches || ret.Years <= base.Years {
		t.Fatalf("1%% lower return should take longer: base=%v ret=%v", base.Years, ret.Years)
	}
}

func TestHouseMoveChangesYears(t *testing.T) {
	in := DefaultFIRE()
	in.Housing = "rent"
	base := FIRE(in)
	house := YearMoves(in)[2]
	if house.Kind != "house" {
		t.Fatalf("kind=%s", house.Kind)
	}
	if !house.Reaches || house.Years <= base.Years {
		t.Fatalf("buying should take longer than renting: rent=%v buy=%v", base.Years, house.Years)
	}
}

func TestPayCutLengthensYears(t *testing.T) {
	in := DefaultFIRE()
	base := FIRE(in)
	opts := YearOptions(in)
	if len(opts) != 3 {
		t.Fatalf("options=%d", len(opts))
	}
	cut := opts[0]
	if cut.Kind != "paycut" || !cut.Reaches || cut.Years <= base.Years {
		t.Fatalf("40%% pay cut should take longer: base=%v cut=%v", base.Years, cut.Years)
	}
	pause := opts[2]
	if pause.Kind != "pause" || !pause.Reaches || pause.Years <= base.Years {
		t.Fatalf("two years off should take longer: base=%v pause=%v", base.Years, pause.Years)
	}
}

func TestMixedUSDAddsToIndiaCorpus(t *testing.T) {
	plain := DefaultFIRE()
	mixed := DefaultFIRE()
	mixed.Mixed = true
	mixed.UsdParked = 10_000
	mixed.UsdMonthly = 500
	a := FIRE(plain)
	b := FIRE(mixed)
	wantExtra := 10_000 * DefaultFxINRPerUSD
	if b.StartingCorpus < a.StartingCorpus+wantExtra-1 {
		t.Fatalf("mixed corpus=%v want at least %v", b.StartingCorpus, a.StartingCorpus+wantExtra)
	}
	if !b.ReachesFire || b.Years >= a.Years {
		t.Fatalf("dollar sleeve should shorten years: plain=%v mixed=%v", a.Years, b.Years)
	}
}

func TestMixedIndiaPotsAddToUSD(t *testing.T) {
	plain := FIREInput{
		Age: 30, AnnualExpenses: 60_000, CurrentSavings: 80_000, MonthlySavings: 2_500,
		ExpectedReturn: 8, Inflation: 3, SWR: 4, Housing: "rent", Region: "us",
	}
	mixed := plain
	mixed.Mixed = true
	mixed.EPFNow = 8_400_000 // ₹84L → $10k at 84
	a := FIRE(plain)
	b := FIRE(mixed)
	if b.StartingCorpus < a.StartingCorpus+10_000-1 {
		t.Fatalf("mixed usd corpus=%v want +10000 vs %v", b.StartingCorpus, a.StartingCorpus)
	}
	if !b.ReachesFire || b.Years >= a.Years {
		t.Fatalf("EPF sleeve should shorten years: plain=%v mixed=%v", a.Years, b.Years)
	}
}
