package calc

import (
	"math"
	"testing"
)

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

func TestLifestyleLaterInflatesTwentyYears(t *testing.T) {
	in := DefaultFIRE()
	in.Housing = "own"
	out := FIRE(in)
	want := out.Lifestyle * math.Pow(1+in.Inflation/100, 20)
	if math.Abs(out.LifestyleLater-want) > 1 {
		t.Fatalf("lifestyle later=%v want %v", out.LifestyleLater, want)
	}
	if math.Abs(out.FireNumberLater-out.FireNumber*math.Pow(1+in.Inflation/100, 20)) > 1 {
		t.Fatalf("number later=%v", out.FireNumberLater)
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
	// 25k/mo equity vs 30L target, 0% return, 0% inflation → ~10 years
	in := FIREInput{
		Age:            30,
		AnnualExpenses: 120_000, // 30L target at 4% SWR
		CurrentSavings: 0,
		ForeignMonthly: 25_000,
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

func TestPotsCountTowardStartingCorpus(t *testing.T) {
	in := DefaultFIRE()
	in.CurrentSavings = 0
	in.NPSNow = 1_000_000
	in.EPFNow = 500_000
	in.PPFNow = 200_000
	in.ForeignNow = 100_000
	in.StoppedNow = 200_000
	out := FIRE(in)
	if out.StartingCorpus != 2_000_000 {
		t.Fatalf("starting corpus=%v want 2000000", out.StartingCorpus)
	}
}

func TestStoppedGrowsWithoutContribution(t *testing.T) {
	in := FIREInput{
		Age:            30,
		AnnualExpenses: 120_000,
		StoppedNow:     1_500_000,
		ExpectedReturn: 12,
		Inflation:      0,
		SWR:            4,
		Housing:        "rent",
	}
	out := FIRE(in)
	if !out.ReachesFire {
		t.Fatal("leftover pot should grow to the 30L target")
	}
	if out.MonthlyIn != 0 {
		t.Fatalf("monthly in=%v want 0", out.MonthlyIn)
	}
}

func TestBuyHouseTier1AddsToNumber(t *testing.T) {
	rent := DefaultFIRE()
	rent.Housing = "rent"
	buy := rent
	buy.Housing = "buy"
	buy.CityTier = 1
	a := FIRE(rent)
	b := FIRE(buy)
	if b.HouseAdd != HouseIN1 {
		t.Fatalf("house add=%v want %v", b.HouseAdd, HouseIN1)
	}
	if b.FireNumber != a.Lifestyle+HouseIN1 {
		t.Fatalf("buy number=%v want lifestyle+house %v", b.FireNumber, a.Lifestyle+HouseIN1)
	}
	if b.Years <= a.Years {
		t.Fatalf("buying a house should take longer: rent %v buy %v", a.Years, b.Years)
	}
}

func TestTier3HouseCheaperThanTier1(t *testing.T) {
	if HouseCost(3, "in") >= HouseCost(1, "in") {
		t.Fatalf("tier 3 should be cheaper: %v vs %v", HouseCost(3, "in"), HouseCost(1, "in"))
	}
	in := DefaultFIRE()
	in.Housing = "buy"
	in.CityTier = 1
	t1 := FIRE(in)
	in.CityTier = 3
	t3 := FIRE(in)
	if t3.FireNumber >= t1.FireNumber {
		t.Fatalf("tier 3 number should be smaller: %v vs %v", t3.FireNumber, t1.FireNumber)
	}
}

func TestUSDHouseUsesDollarScale(t *testing.T) {
	in := DefaultFIRE()
	in.Region = "us"
	in.Housing = "buy"
	in.CityTier = 1
	in.AnnualExpenses = 60_000
	out := FIRE(in)
	if out.HouseAdd != HouseUS1 {
		t.Fatalf("usd house=%v want %v", out.HouseAdd, HouseUS1)
	}
}

func TestUSDRetirementGrowsAtExpectedReturn(t *testing.T) {
	base := FIREInput{
		Age:            30,
		AnnualExpenses: 40_000,
		NPSMonthly:     2_000,
		ExpectedReturn: 12,
		Inflation:      0,
		SWR:            4,
		Housing:        "rent",
	}
	in := base
	in.Region = "in"
	us := base
	us.Region = "us"
	a := FIRE(in)
	b := FIRE(us)
	if !a.ReachesFire || !b.ReachesFire {
		t.Fatalf("both should reach: in=%v us=%v", a.ReachesFire, b.ReachesFire)
	}
	if b.Years >= a.Years {
		t.Fatalf("USD retirement at 12%% should beat NPS 9%%: in %v us %v", a.Years, b.Years)
	}
}

func TestLiquidGrowsAtParkingNotEquity(t *testing.T) {
	liq := FIREInput{
		Age:            30,
		AnnualExpenses: 120_000,
		CurrentSavings: 1_500_000,
		ExpectedReturn: 12,
		Inflation:      0,
		SWR:            4,
		Housing:        "rent",
	}
	eq := liq
	eq.CurrentSavings = 0
	eq.StoppedNow = 1_500_000
	a := FIRE(liq)
	b := FIRE(eq)
	if !a.ReachesFire || !b.ReachesFire {
		t.Fatalf("both should reach: liquid=%v equity=%v", a.ReachesFire, b.ReachesFire)
	}
	if b.Years >= a.Years {
		t.Fatalf("leftover SIPs at 12%% should beat liquid parking: liquid %v leftover %v", a.Years, b.Years)
	}
}

func TestGoldCountsAndBeatsParked(t *testing.T) {
	in := DefaultFIRE()
	in.CurrentSavings = 0
	in.GoldNow = 1_000_000
	out := FIRE(in)
	if out.StartingCorpus != 1_000_000 {
		t.Fatalf("gold starting corpus=%v want 1000000", out.StartingCorpus)
	}
	parked := FIREInput{
		Age:            30,
		AnnualExpenses: 120_000,
		CurrentSavings: 1_500_000,
		ExpectedReturn: 12,
		Inflation:      0,
		SWR:            4,
		Housing:        "rent",
	}
	gold := parked
	gold.CurrentSavings = 0
	gold.GoldNow = 1_500_000
	a := FIRE(parked)
	b := FIRE(gold)
	if !a.ReachesFire || !b.ReachesFire {
		t.Fatalf("both should reach: parked=%v gold=%v", a.ReachesFire, b.ReachesFire)
	}
	if b.Years >= a.Years {
		t.Fatalf("gold at 8%% should beat parked 6%%: parked %v gold %v", a.Years, b.Years)
	}
}

func TestJewelleryIsAssetNotSpendable(t *testing.T) {
	base := FIREInput{
		Age:            30,
		AnnualExpenses: 120_000,
		ExpectedReturn: 8,
		Inflation:      0,
		SWR:            4,
		Housing:        "rent",
	}
	jew := base
	jew.JewelleryNow = 1_500_000
	gold := base
	gold.GoldNow = 1_500_000
	a := FIRE(base)
	j := FIRE(jew)
	g := FIRE(gold)
	if j.StartingCorpus != 1_500_000 {
		t.Fatalf("jewellery should count as an asset: starting=%v", j.StartingCorpus)
	}
	if j.Jewellery != 1_500_000 {
		t.Fatalf("jewellery field=%v want 1500000", j.Jewellery)
	}
	if j.ReachesFire != a.ReachesFire || j.Years != a.Years {
		t.Fatalf("jewellery you keep should not change years-to-FIRE: base years=%v jew years=%v", a.Years, j.Years)
	}
	growing := DefaultFIRE()
	growing.JewelleryNow = 1_000_000
	grown := FIRE(growing)
	if grown.JewelleryLater <= grown.Jewellery {
		t.Fatalf("jewellery should grow in net worth: now=%v later=%v", grown.Jewellery, grown.JewelleryLater)
	}
	plain := DefaultFIRE()
	plainYears := FIRE(plain)
	if grown.Years != plainYears.Years {
		t.Fatalf("growing jewellery must not change FIRE years: %v vs %v", grown.Years, plainYears.Years)
	}
	if !g.ReachesFire {
		t.Fatal("sellable gold should help reach FIRE")
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

func TestStepUpShortensYears(t *testing.T) {
	flat := DefaultFIRE()
	flat.StepUp = 0
	a := FIRE(flat)
	up := DefaultFIRE()
	up.StepUp = 10
	b := FIRE(up)
	if !a.ReachesFire || !b.ReachesFire {
		t.Fatal("both should reach")
	}
	if b.Years >= a.Years {
		t.Fatalf("10%% step-up should be faster: %v vs %v", b.Years, a.Years)
	}
}

func TestFireChartSplitsPots(t *testing.T) {
	in := DefaultFIRE()
	in.StoppedNow = 400_000
	in.GoldNow = 100_000
	in.NPSNow = 200_000
	out := FIRE(in)
	if len(out.Chart) < 2 {
		t.Fatal("expected chart points")
	}
	p0 := out.Chart[0]
	if p0.Parked != in.CurrentSavings {
		t.Fatalf("year 0 parked=%v want %v", p0.Parked, in.CurrentSavings)
	}
	if p0.Invested != in.StoppedNow {
		t.Fatalf("year 0 invested=%v want %v", p0.Invested, in.StoppedNow)
	}
	if p0.Gold != in.GoldNow || p0.NPS != in.NPSNow {
		t.Fatalf("year 0 gold=%v nps=%v", p0.Gold, p0.NPS)
	}
	last := out.Chart[len(out.Chart)-1]
	if last.Invested <= p0.Invested {
		t.Fatalf("invested SIP pot should grow: start=%v end=%v", p0.Invested, last.Invested)
	}
	var y20 *FIREPoint
	for i := range out.Chart {
		if out.Chart[i].Year == 20 {
			y20 = &out.Chart[i]
			break
		}
	}
	if y20 == nil {
		t.Fatalf("expected a year-20 snapshot, last year=%v", last.Year)
	}
}
