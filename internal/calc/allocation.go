package calc

import (
	"math"
	"strconv"
)

const (
	CatLarge = "Large-cap index"
	CatMid   = "Mid-cap index"
	CatSmall = "Small-cap index"
	CatDebt  = "Debt / government securities"
	CatGold  = "Gold"
	CatIntl  = "International equity index"
	CatCash  = "Cash / liquid"
)

var AllowedCategories = []string{
	CatLarge, CatMid, CatSmall, CatDebt, CatGold, CatIntl, CatCash,
}

type Sleeve struct {
	Category string  `json:"category"`
	Percent  float64 `json:"percent"`
}

type AllocationResult struct {
	Sleeves   []Sleeve `json:"sleeves"`
	Rationale string   `json:"rationale"`
	Source    string   `json:"source"`
}

type sleeves struct {
	large, mid, small, debt, gold, intl, cash float64
}

func HeuristicAllocation(age, horizon int, risk string) AllocationResult {
	var s sleeves
	switch risk {
	case "conservative":
		s = sleeves{30, 5, 0, 40, 10, 5, 10}
	case "aggressive":
		s = sleeves{40, 20, 10, 10, 5, 15, 0}
	default:
		risk = "moderate"
		s = sleeves{40, 15, 5, 20, 8, 10, 2}
	}

	if horizon <= 3 {
		amt := 20.0
		if risk == "aggressive" {
			amt = 12
		}
		taken := takeFromEquity(&s, amt)
		s.debt += taken * 0.70
		s.cash += taken * 0.30
	} else if horizon <= 7 {
		taken := takeFromEquity(&s, 10)
		s.debt += taken
	}

	if age >= 50 {
		taken := takeFromEquity(&s, 8)
		s.debt += taken * 0.60
		s.gold += taken * 0.40
		if s.small > 0 {
			s.debt += s.small
			s.small = 0
		}
	}

	list := normalizeSleeves(s)
	return AllocationResult{
		Sleeves:   list,
		Rationale: heuristicRationale(age, horizon, risk),
		Source:    "On-device heuristic",
	}
}

func heuristicRationale(age, horizon int, risk string) string {
	return "Category mix from a simple offline heuristic for age " +
		strconv.Itoa(age) + " and a " + strconv.Itoa(horizon) + "-year horizon: a " + risk +
		" base, then shorter horizons shift equity into debt and cash, and age 50+ trims small-cap toward debt and gold. " +
		"Sleeves are asset classes only — not products. This is educational arithmetic, not a recommendation."
}

func takeFromEquity(s *sleeves, amt float64) float64 {
	taken := 0.0
	order := []*float64{&s.small, &s.mid, &s.large}
	for _, ptr := range order {
		need := amt - taken
		if need <= 0 {
			break
		}
		if *ptr >= need {
			*ptr -= need
			taken += need
		} else {
			taken += *ptr
			*ptr = 0
		}
	}
	return taken
}

func normalizeSleeves(s sleeves) []Sleeve {
	raw := []Sleeve{
		{Category: CatLarge, Percent: s.large},
		{Category: CatMid, Percent: s.mid},
		{Category: CatSmall, Percent: s.small},
		{Category: CatDebt, Percent: s.debt},
		{Category: CatGold, Percent: s.gold},
		{Category: CatIntl, Percent: s.intl},
		{Category: CatCash, Percent: s.cash},
	}
	return NormalizePercents(raw)
}

// NormalizePercents scales sleeves to 100 and rounds to one decimal.
func NormalizePercents(in []Sleeve) []Sleeve {
	sum := 0.0
	for _, s := range in {
		if s.Percent > 0 {
			sum += s.Percent
		}
	}
	if sum <= 0 {
		return []Sleeve{
			{Category: CatLarge, Percent: 40},
			{Category: CatMid, Percent: 15},
			{Category: CatSmall, Percent: 5},
			{Category: CatDebt, Percent: 20},
			{Category: CatGold, Percent: 8},
			{Category: CatIntl, Percent: 10},
			{Category: CatCash, Percent: 2},
		}
	}
	out := make([]Sleeve, 0, len(in))
	for _, s := range in {
		if s.Percent < 0 {
			continue
		}
		out = append(out, Sleeve{Category: s.Category, Percent: s.Percent / sum * 100})
	}
	// Round to 1 decimal and fix remainder on the largest sleeve.
	roundedSum := 0.0
	maxI := 0
	for i := range out {
		out[i].Percent = math.Round(out[i].Percent*10) / 10
		roundedSum += out[i].Percent
		if out[i].Percent > out[maxI].Percent {
			maxI = i
		}
	}
	diff := math.Round((100-roundedSum)*10) / 10
	out[maxI].Percent += diff
	if out[maxI].Percent < 0 {
		out[maxI].Percent = 0
	}
	return out
}

func SleevesSum(ss []Sleeve) float64 {
	t := 0.0
	for _, s := range ss {
		t += s.Percent
	}
	return t
}
