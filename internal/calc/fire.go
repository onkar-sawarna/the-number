package calc

import "math"

type FIREInput struct {
	Age            int     `json:"age"`
	AnnualExpenses float64 `json:"annual_expenses"`
	CurrentSavings float64 `json:"current_savings"`
	MonthlySavings float64 `json:"monthly_savings"`
	ExpectedReturn float64 `json:"expected_return"`
	Inflation      float64 `json:"inflation"`
	SWR            float64 `json:"swr"`
}

type FIREPoint struct {
	Year   int     `json:"year"`
	Age    int     `json:"age"`
	Corpus float64 `json:"corpus"`
	Target float64 `json:"target"`
}

type FIREOutput struct {
	FireNumber  float64     `json:"fire_number"`
	Lean        float64     `json:"lean"`
	Regular     float64 `json:"regular"`
	Fat         float64     `json:"fat"`
	Years       float64     `json:"years"`
	ReachesFire bool        `json:"reaches_fire"`
	FIAge       int         `json:"fi_age"`
	Chart       []FIREPoint `json:"chart"`
}

func DefaultFIRE() FIREInput {
	return FIREInput{
		Age:            30,
		AnnualExpenses: 1_200_000,
		CurrentSavings: 1_500_000,
		MonthlySavings: 50_000,
		ExpectedReturn: 11,
		Inflation:      6,
		SWR:            4,
	}
}

func FIRENumber(annualExpenses, swr float64) float64 {
	if swr <= 0 {
		return 0
	}
	return annualExpenses / (swr / 100.0)
}

func FIRE(in FIREInput) FIREOutput {
	swr := in.SWR
	if swr <= 0 {
		swr = 4
	}
	fireNum := FIRENumber(in.AnnualExpenses, swr)
	out := FIREOutput{
		FireNumber: fireNum,
		Lean:       fireNum * 0.5,
		Regular:    fireNum,
		Fat:        fireNum * 2,
		FIAge:      in.Age,
	}

	if in.CurrentSavings >= fireNum {
		out.Years = 0
		out.ReachesFire = true
		out.Chart = fireChart(in, swr, 0, true)
		return out
	}

	rm := monthlyEffective(in.ExpectedReturn)
	im := monthlyEffective(in.Inflation)
	corpus := in.CurrentSavings
	expenses := in.AnnualExpenses
	const maxMonths = 80 * 12
	for m := 1; m <= maxMonths; m++ {
		corpus = corpus*(1+rm) + in.MonthlySavings
		expenses = expenses * (1 + im)
		target := FIRENumber(expenses, swr)
		if corpus >= target {
			out.ReachesFire = true
			out.Years = float64(m) / 12.0
			out.FIAge = in.Age + int(math.Round(out.Years))
			out.Chart = fireChart(in, swr, out.Years, true)
			return out
		}
	}
	out.ReachesFire = false
	out.Years = 0
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
	rm := monthlyEffective(in.ExpectedReturn)
	inf := in.Inflation / 100.0
	corpus := in.CurrentSavings
	expenses := in.AnnualExpenses
	pts := make([]FIREPoint, 0, horizon+1)
	for y := 0; y <= horizon; y++ {
		pts = append(pts, FIREPoint{
			Year:   y,
			Age:    in.Age + y,
			Corpus: corpus,
			Target: FIRENumber(expenses, swr),
		})
		for m := 0; m < 12; m++ {
			corpus = corpus*(1+rm) + in.MonthlySavings
		}
		expenses *= (1 + inf)
	}
	return pts
}
