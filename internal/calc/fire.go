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

	NPSNow        float64 `json:"nps_now"`
	NPSMonthly    float64 `json:"nps_monthly"`
	PPFNow        float64 `json:"ppf_now"`
	PPFMonthly    float64 `json:"ppf_monthly"`
	EPFNow        float64 `json:"epf_now"`
	EPFMonthly    float64 `json:"epf_monthly"`
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
	FireNumber     float64     `json:"fire_number"`
	Lifestyle      float64     `json:"lifestyle"`
	HouseAdd       float64     `json:"house_add"`
	StartingCorpus float64     `json:"starting_corpus"`
	Jewellery      float64     `json:"jewellery"`
	JewelleryLater float64     `json:"jewellery_later"`
	MonthlyIn      float64     `json:"monthly_in"`
	Lean           float64     `json:"lean"`
	Regular        float64     `json:"regular"`
	Fat            float64     `json:"fat"`
	Years          float64     `json:"years"`
	ReachesFire    bool        `json:"reaches_fire"`
	FIAge          int         `json:"fi_age"`
	Chart          []FIREPoint `json:"chart"`
}

func DefaultFIRE() FIREInput {
	return FIREInput{
		Age:            30,
		AnnualExpenses: 1_200_000,
		CurrentSavings: 1_500_000,
		MonthlySavings: 50_000,
		ExpectedReturn: 12,
		Inflation:      6,
		SWR:            3.5,
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
		general: nz(in.CurrentSavings),
		nps:     nz(in.NPSNow),
		ppf:     nz(in.PPFNow),
		epf:     nz(in.EPFNow),
		foreign: nz(in.ForeignNow),
		stopped: nz(in.StoppedNow),
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

func stepPots(p firePots, in FIREInput, month int, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold float64) firePots {
	p.general = p.general * (1 + rmLiq)
	p.nps = p.nps*(1+rmNPS) + steppedMonthly(in.NPSMonthly, in.StepUp, month)
	p.ppf = p.ppf*(1+rmPPF) + ppfAt(in, month)
	p.epf = p.epf*(1+rmEPF) + steppedMonthly(in.EPFMonthly, in.StepUp, month)
	p.foreign = p.foreign*(1+rmEq) + steppedMonthly(in.ForeignMonthly, in.StepUp, month)
	p.stopped = p.stopped*(1+rmEq) + steppedMonthly(in.MonthlySavings, in.StepUp, month)
	p.gold = p.gold*(1+rmGold) + steppedMonthly(in.GoldMonthly, in.StepUp, month)
	p.jewellery = p.jewellery * (1 + rmGold)
	return p
}

func FIRE(in FIREInput) FIREOutput {
	swr := in.SWR
	if swr <= 0 {
		swr = 4
	}
	lifestyle := FIRENumber(in.AnnualExpenses, swr)
	house := HouseAdd(in)
	fireNum := lifestyle + house
	pots := startPots(in)
	out := FIREOutput{
		FireNumber:     fireNum,
		Lifestyle:      lifestyle,
		HouseAdd:       house,
		StartingCorpus: pots.total(),
		Jewellery:      pots.jewellery,
		JewelleryLater: pots.jewellery,
		MonthlyIn:      monthlyIn(in),
		Lean:           lifestyle*0.5 + house,
		Regular:        fireNum,
		Fat:            lifestyle*2 + house,
		FIAge:          in.Age,
	}

	if pots.spendable() >= fireNum {
		out.Years = 0
		out.ReachesFire = true
		out.Chart = fireChart(in, swr, 0, true)
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
	const maxMonths = 80 * 12
	for m := 1; m <= maxMonths; m++ {
		pots = stepPots(pots, in, m, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold)
		expenses = expenses * (1 + im)
		houseNow = houseNow * (1 + im)
		target := FIRENumber(expenses, swr) + houseNow
		if pots.spendable() >= target {
			out.ReachesFire = true
			out.Years = float64(m) / 12.0
			out.FIAge = in.Age + int(math.Round(out.Years))
			out.JewelleryLater = pots.jewellery
			out.Chart = fireChart(in, swr, out.Years, true)
			return out
		}
	}
	out.ReachesFire = false
	out.Years = 0
	out.JewelleryLater = pots.jewellery
	out.Chart = fireChart(in, swr, 0, false)
	return out
}

func fireChart(in FIREInput, swr, years float64, reaches bool) []FIREPoint {
	horizon := 40
	if reaches {
		horizon = int(math.Ceil(years)) + 5
		if horizon < 15 {
			horizon = 15
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
