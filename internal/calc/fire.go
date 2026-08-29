package calc

import "math"

// Educational sleeve rates for Indian pots. Parked cash uses RateLiquid.
// Gold funds use RateGold. Jewellery also grows at RateGold as a kept net-worth asset; it is not spendable for FIRE.
const (
	RateNPS    = 9.0  // blended NPS
	RatePPF    = 7.1  // current PPF-like rate
	RateEPF    = 8.25 // EPF-like rate
	RateLiquid = 6.0  // savings / FD-like parking, not equity SIPs
	RateGold   = 8.0  // long-run gold-like, not a forecast
)

// Indicative modest house (today’s money). Not a quote.
const (
	HouseIN1 = 20_000_000.0
	HouseIN2 = 9_000_000.0
	HouseIN3 = 4_500_000.0
	HouseUS1 = 800_000.0
	HouseUS2 = 400_000.0
	HouseUS3 = 220_000.0
)

type FIREInput struct {
	Age            int     `json:"age"`
	AnnualExpenses float64 `json:"annual_expenses"`
	CurrentSavings float64 `json:"current_savings"`
	MonthlySavings float64 `json:"monthly_savings"`
	ExpectedReturn float64 `json:"expected_return"`
	Inflation      float64 `json:"inflation"`
	SWR            float64 `json:"swr"`

	NPSNow         float64 `json:"nps_now"`
	NPSMonthly     float64 `json:"nps_monthly"`
	PPFNow         float64 `json:"ppf_now"`
	PPFMonthly     float64 `json:"ppf_monthly"`
	EPFNow         float64 `json:"epf_now"`
	EPFMonthly     float64 `json:"epf_monthly"`
	ForeignNow     float64 `json:"foreign_now"`
	ForeignMonthly float64 `json:"foreign_monthly"`
	StoppedNow     float64 `json:"stopped_now"`
	GoldNow        float64 `json:"gold_now"`
	GoldMonthly    float64 `json:"gold_monthly"`
	JewelleryNow   float64 `json:"jewellery_now"`
	StepUp         float64 `json:"step_up"` // yearly % rise in monthly SIPs

	CityTier int    `json:"city_tier"` // 1, 2, or 3
	Housing  string `json:"housing"`   // own | rent | buy
	Region   string `json:"region"`    // in | us

	ContribScale    float64 `json:"contrib_scale"`     // 0 = 1 (full). 0.6 = 40% pay cut.
	PauseMonths     int     `json:"pause_months"`      // first N months, deposits are 0
	StopAfterMonths int     `json:"stop_after_months"` // 0 = keep contributing. N > 0 = deposits stop after month N.
	SkipChart       bool    `json:"skip_chart"`        // Coast trials skip the yearly chart.
}

type firePots struct {
	general, nps, ppf, epf, foreign, stopped, gold, jewellery float64
}

func (p firePots) spendable() float64 {
	return p.general + p.nps + p.ppf + p.epf + p.foreign + p.stopped + p.gold
}

func (p firePots) total() float64 {
	return p.spendable() + p.jewellery
}

type FIREPoint struct {
	Year      int     `json:"year"`
	Age       int     `json:"age"`
	Corpus    float64 `json:"corpus"`
	NetWorth  float64 `json:"net_worth"`
	Target    float64 `json:"target"`
	Parked    float64 `json:"parked"`
	NPS       float64 `json:"nps"`
	PPF       float64 `json:"ppf"`
	EPF       float64 `json:"epf"`
	Foreign   float64 `json:"foreign"`
	Invested  float64 `json:"invested"`
	Gold      float64 `json:"gold"`
	Jewellery float64 `json:"jewellery"`
}

type FIREOutput struct {
	FireNumber      float64     `json:"fire_number"`
	FireNumberLater float64     `json:"fire_number_later"`
	Lifestyle       float64     `json:"lifestyle"`
	LifestyleLater  float64     `json:"lifestyle_later"`
	HouseAdd        float64     `json:"house_add"`
	StartingCorpus  float64     `json:"starting_corpus"`
	SpendableNow    float64     `json:"spendable_now"`
	StillNeed       float64     `json:"still_need"`
	Jewellery       float64     `json:"jewellery"`
	JewelleryLater  float64     `json:"jewellery_later"`
	MonthlyIn       float64     `json:"monthly_in"`
	Lean            float64     `json:"lean"`
	Regular         float64     `json:"regular"`
	Fat             float64     `json:"fat"`
	Years           float64     `json:"years"`
	ReachesFire     bool        `json:"reaches_fire"`
	FIAge           int         `json:"fi_age"`
	CrossingYears   float64     `json:"crossing_years"`
	ReachesCrossing bool        `json:"reaches_crossing"`
	CrossingAge     int         `json:"crossing_age"`
	Chart           []FIREPoint `json:"chart"`
}

func DefaultFIRE() FIREInput {
	return FIREInput{
		Age:            30,
		AnnualExpenses: 1_200_000,
		CurrentSavings: 1_500_000,
		MonthlySavings: 50_000,
		ExpectedReturn: 12,
		Inflation:      6,
		SWR:            4,
		StepUp:         10,
		CityTier:       1,
		Housing:        "rent",
		Region:         "in",
	}
}

func IsUSD(region string) bool {
	return region == "us" || region == "usd"
}

func HouseCost(tier int, region string) float64 {
	if IsUSD(region) {
		switch tier {
		case 3:
			return HouseUS3
		case 2:
			return HouseUS2
		default:
			return HouseUS1
		}
	}
	switch tier {
	case 3:
		return HouseIN3
	case 2:
		return HouseIN2
	default:
		return HouseIN1
	}
}

func HouseAdd(in FIREInput) float64 {
	if in.Housing != "buy" {
		return 0
	}
	return HouseCost(in.CityTier, in.Region)
}

func FIRENumber(annualExpenses, swr float64) float64 {
	if swr <= 0 {
		return 0
	}
	return annualExpenses / (swr / 100.0)
}

func startPots(in FIREInput) firePots {
	return firePots{
		general:   nz(in.CurrentSavings),
		nps:       nz(in.NPSNow),
		ppf:       nz(in.PPFNow),
		epf:       nz(in.EPFNow),
		foreign:   nz(in.ForeignNow),
		stopped:   nz(in.StoppedNow),
		gold:      nz(in.GoldNow),
		jewellery: nz(in.JewelleryNow),
	}
}

func nz(v float64) float64 {
	if v < 0 {
		return 0
	}
	return v
}

func monthlyIn(in FIREInput) float64 {
	return nz(in.MonthlySavings) + nz(in.NPSMonthly) + nz(in.PPFMonthly) + nz(in.EPFMonthly) + nz(in.ForeignMonthly) + nz(in.GoldMonthly)
}

func steppedMonthly(base, stepUpPct float64, month int) float64 {
	base = nz(base)
	if base == 0 || stepUpPct <= 0 || month <= 12 {
		return base
	}
	years := (month - 1) / 12
	return base * math.Pow(1+stepUpPct/100.0, float64(years))
}

func ppfAt(in FIREInput, month int) float64 {
	v := steppedMonthly(in.PPFMonthly, in.StepUp, month)
	if !IsUSD(in.Region) && v > 12_500 {
		return 12_500
	}
	return v
}

func depositScale(in FIREInput, month int) float64 {
	if in.StopAfterMonths < 0 {
		return 0
	}
	if in.StopAfterMonths > 0 && month > in.StopAfterMonths {
		return 0
	}
	if in.PauseMonths > 0 && month <= in.PauseMonths {
		return 0
	}
	if in.ContribScale > 0 && in.ContribScale != 1 {
		return in.ContribScale
	}
	return 1
}

func addDeposit(base float64, in FIREInput, month int) float64 {
	return steppedMonthly(base, in.StepUp, month) * depositScale(in, month)
}

func stepPots(p firePots, in FIREInput, month int, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold float64) firePots {
	s := depositScale(in, month)
	p.general = p.general * (1 + rmLiq)
	p.nps = p.nps*(1+rmNPS) + addDeposit(in.NPSMonthly, in, month)
	p.ppf = p.ppf*(1+rmPPF) + ppfAt(in, month)*s
	p.epf = p.epf*(1+rmEPF) + addDeposit(in.EPFMonthly, in, month)
	p.foreign = p.foreign*(1+rmEq) + addDeposit(in.ForeignMonthly, in, month)
	p.stopped = p.stopped*(1+rmEq) + addDeposit(in.MonthlySavings, in, month)
	p.gold = p.gold*(1+rmGold) + addDeposit(in.GoldMonthly, in, month)
	p.jewellery = p.jewellery * (1 + rmGold)
	return p
}

func contribAt(in FIREInput, month int) float64 {
	s := depositScale(in, month)
	return addDeposit(in.NPSMonthly, in, month) +
		ppfAt(in, month)*s +
		addDeposit(in.EPFMonthly, in, month) +
		addDeposit(in.ForeignMonthly, in, month) +
		addDeposit(in.MonthlySavings, in, month) +
		addDeposit(in.GoldMonthly, in, month)
}

func FIRE(in FIREInput) FIREOutput {
	swr := in.SWR
	if swr <= 0 {
		swr = 4
	}
	lifestyle := FIRENumber(in.AnnualExpenses, swr)
	house := HouseAdd(in)
	fireNum := lifestyle + house
	grow := math.Pow(1+in.Inflation/100.0, 20)
	pots := startPots(in)
	spendable := pots.spendable()
	still := fireNum - spendable
	if still < 0 {
		still = 0
	}
	out := FIREOutput{
		FireNumber:      fireNum,
		FireNumberLater: fireNum * grow,
		Lifestyle:       lifestyle,
		LifestyleLater:  lifestyle * grow,
		HouseAdd:        house,
		StartingCorpus:  pots.total(),
		SpendableNow:    spendable,
		StillNeed:       still,
		Jewellery:       pots.jewellery,
		JewelleryLater:  pots.jewellery,
		MonthlyIn:       monthlyIn(in),
		Lean:            nz(in.AnnualExpenses)*20 + house,
		Regular:         fireNum,
		Fat:             nz(in.AnnualExpenses)*50 + house,
		FIAge:           in.Age,
		CrossingAge:     in.Age,
	}

	fireSet := false
	if pots.spendable() >= fireNum {
		out.Years = 0
		out.ReachesFire = true
		out.JewelleryLater = pots.jewellery
		fireSet = true
	}
	if monthlyIn(in) <= 0 {
		if pots.spendable() > 0 {
			out.ReachesCrossing = true
			out.CrossingYears = 0
			out.CrossingAge = in.Age
		}
	}
	if fireSet && out.ReachesCrossing {
		out.Chart = maybeChart(in, swr, out.Years, out.ReachesFire)
		return out
	}

	rmLiq := monthlyEffective(RateLiquid)
	rmEq := monthlyEffective(in.ExpectedReturn)
	rmNPS := monthlyEffective(RateNPS)
	if IsUSD(in.Region) {
		rmNPS = rmEq
	}
	rmPPF := monthlyEffective(RatePPF)
	rmEPF := monthlyEffective(RateEPF)
	rmGold := monthlyEffective(RateGold)
	im := monthlyEffective(in.Inflation)
	expenses := in.AnnualExpenses
	houseNow := house
	startYear := pots.spendable()
	contribYear := 0.0
	const maxMonths = 80 * 12
	for m := 1; m <= maxMonths; m++ {
		contribYear += contribAt(in, m)
		pots = stepPots(pots, in, m, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold)
		expenses = expenses * (1 + im)
		houseNow = houseNow * (1 + im)
		if !fireSet {
			target := FIRENumber(expenses, swr) + houseNow
			if pots.spendable() >= target {
				fireSet = true
				out.ReachesFire = true
				out.Years = float64(m) / 12.0
				out.FIAge = in.Age + int(math.Round(out.Years))
				out.JewelleryLater = pots.jewellery
			}
		}
		if m%12 == 0 && !out.ReachesCrossing {
			end := pots.spendable()
			growth := end - startYear - contribYear
			if contribYear > 0 && growth+1e-6 >= contribYear {
				out.ReachesCrossing = true
				out.CrossingYears = float64(m) / 12.0
				out.CrossingAge = in.Age + int(math.Round(out.CrossingYears))
			}
			startYear = end
			contribYear = 0
		}
		if fireSet && out.ReachesCrossing {
			break
		}
	}
	if !fireSet {
		out.ReachesFire = false
		out.Years = 0
		out.JewelleryLater = pots.jewellery
	}
	out.Chart = maybeChart(in, swr, out.Years, out.ReachesFire)
	return out
}

func maybeChart(in FIREInput, swr, years float64, reaches bool) []FIREPoint {
	if in.SkipChart {
		return nil
	}
	return fireChart(in, swr, years, reaches)
}

// CoastOutput is the earliest you can stop SIPs and still reach FIRE.
type CoastOutput struct {
	Reaches   bool    `json:"reaches"`
	Already   bool    `json:"already"`
	Years     float64 `json:"years"`
	Age       int     `json:"age"`
	LandYears float64 `json:"landYears"`
	LandAge   int     `json:"landAge"`
	UntilFire bool    `json:"untilFire"`
}

func zeroMonthlies(in FIREInput) FIREInput {
	in.MonthlySavings = 0
	in.NPSMonthly = 0
	in.PPFMonthly = 0
	in.EPFMonthly = 0
	in.ForeignMonthly = 0
	in.GoldMonthly = 0
	in.StopAfterMonths = 0
	in.SkipChart = true
	return in
}

// Coast is the fewest months of contributing after which you can stop SIPs
// and still reach FIRE within 80 years. Already means you can stop today.
func Coast(in FIREInput) CoastOutput {
	fullIn := in
	fullIn.SkipChart = true
	full := FIRE(fullIn)
	if !full.ReachesFire {
		return CoastOutput{}
	}
	if full.Years == 0 {
		return CoastOutput{Reaches: true, Already: true, Age: in.Age, LandAge: in.Age}
	}
	deadline := coastDeadlineAge(full.FIAge)
	now := FIRE(zeroMonthlies(in))
	if coastLands(now, deadline) {
		return CoastOutput{
			Reaches:   true,
			Already:   true,
			Age:       in.Age,
			LandYears: now.Years,
			LandAge:   now.FIAge,
		}
	}
	if monthlyIn(in) <= 0 {
		return CoastOutput{
			Reaches:   true,
			Already:   true,
			Age:       in.Age,
			LandYears: full.Years,
			LandAge:   full.FIAge,
		}
	}
	hi := int(math.Ceil(full.Years * 12))
	if hi < 1 {
		hi = 1
	}
	lo := 1
	best := hi
	bestOut := full
	for lo <= hi {
		mid := (lo + hi) / 2
		trial := in
		trial.StopAfterMonths = mid
		trial.SkipChart = true
		out := FIRE(trial)
		if coastLands(out, deadline) {
			best = mid
			bestOut = out
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	years := float64(best) / 12.0
	return CoastOutput{
		Reaches:   true,
		Years:     years,
		Age:       in.Age + int(math.Round(years)),
		LandYears: bestOut.Years,
		LandAge:   bestOut.FIAge,
		UntilFire: best >= int(math.Ceil(full.Years*12)),
	}
}

func coastDeadlineAge(fiAge int) int {
	if fiAge > 60 {
		return fiAge
	}
	return 60
}

func coastLands(out FIREOutput, deadlineAge int) bool {
	return out.ReachesFire && out.FIAge <= deadlineAge
}

func fireChart(in FIREInput, swr, years float64, reaches bool) []FIREPoint {
	horizon := 40
	if reaches {
		horizon = int(math.Ceil(years)) + 5
		if horizon < 20 {
			horizon = 20
		}
		if horizon > 50 {
			horizon = 50
		}
	}
	rmLiq := monthlyEffective(RateLiquid)
	rmEq := monthlyEffective(in.ExpectedReturn)
	rmNPS := monthlyEffective(RateNPS)
	if IsUSD(in.Region) {
		rmNPS = rmEq
	}
	rmPPF := monthlyEffective(RatePPF)
	rmEPF := monthlyEffective(RateEPF)
	rmGold := monthlyEffective(RateGold)
	inf := in.Inflation / 100.0
	pots := startPots(in)
	expenses := in.AnnualExpenses
	house := HouseAdd(in)
	pts := make([]FIREPoint, 0, horizon+1)
	for y := 0; y <= horizon; y++ {
		pts = append(pts, FIREPoint{
			Year:      y,
			Age:       in.Age + y,
			Corpus:    pots.spendable(),
			NetWorth:  pots.total(),
			Target:    FIRENumber(expenses, swr) + house,
			Parked:    pots.general,
			NPS:       pots.nps,
			PPF:       pots.ppf,
			EPF:       pots.epf,
			Foreign:   pots.foreign,
			Invested:  pots.stopped,
			Gold:      pots.gold,
			Jewellery: pots.jewellery,
		})
		for m := 0; m < 12; m++ {
			pots = stepPots(pots, in, y*12+m+1, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold)
		}
		expenses *= (1 + inf)
		house *= (1 + inf)
	}
	return pts
}
