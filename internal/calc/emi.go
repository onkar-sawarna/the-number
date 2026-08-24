package calc

import "math"

type EMIInput struct {
	Principal    float64 `json:"principal"`
	AnnualRate   float64 `json:"annual_rate"`
	Years        int     `json:"years"`
	ExtraMonthly float64 `json:"extra_monthly"`
	Lump         float64 `json:"lump"`
}

type EMIPoint struct {
	Year           int     `json:"year"`
	Balance        float64 `json:"balance"`
	PrepaidBalance float64 `json:"prepaid_balance"`
}

type EMIOutput struct {
	EMI             float64    `json:"emi"`
	Months          int        `json:"months"`
	TotalInterest   float64    `json:"total_interest"`
	TotalPayment    float64    `json:"total_payment"`
	PrepaidMonths   int        `json:"prepaid_months"`
	PrepaidInterest float64    `json:"prepaid_interest"`
	InterestSaved   float64    `json:"interest_saved"`
	MonthsSaved     int        `json:"months_saved"`
	Chart           []EMIPoint `json:"chart"`
}

type loanSim struct {
	opening  float64
	months   int
	interest float64
	monthly  []float64 // balance after month 1..n
}

func DefaultEMI() EMIInput {
	return EMIInput{
		Principal:    2_500_000,
		AnnualRate:   8.5,
		Years:        20,
		ExtraMonthly: 0,
		Lump:         0,
	}
}

func EMIAmount(principal, annualRate float64, years int) float64 {
	n := years * 12
	if n <= 0 {
		return 0
	}
	r := annualRate / 12.0 / 100.0
	if nearlyZero(r) {
		return principal / float64(n)
	}
	pow := math.Pow(1+r, float64(n))
	return principal * r * pow / (pow - 1)
}

func EMI(in EMIInput) EMIOutput {
	if in.Years < 1 {
		in.Years = 1
	}
	n := in.Years * 12
	r := in.AnnualRate / 12.0 / 100.0
	emi := EMIAmount(in.Principal, in.AnnualRate, in.Years)
	base := simulateLoan(in.Principal, r, emi, n)

	prepaidPrincipal := in.Principal - in.Lump
	if prepaidPrincipal < 0 {
		prepaidPrincipal = 0
	}
	prepaid := simulateLoan(prepaidPrincipal, r, emi+in.ExtraMonthly, n)

	out := EMIOutput{
		EMI:             emi,
		Months:          base.months,
		TotalInterest:   base.interest,
		TotalPayment:    in.Principal + base.interest,
		PrepaidMonths:   prepaid.months,
		PrepaidInterest: prepaid.interest,
		InterestSaved:   base.interest - prepaid.interest,
		MonthsSaved:     base.months - prepaid.months,
	}
	if out.InterestSaved < 0 {
		out.InterestSaved = 0
	}
	if out.MonthsSaved < 0 {
		out.MonthsSaved = 0
	}

	out.Chart = make([]EMIPoint, 0, in.Years+1)
	for y := 0; y <= in.Years; y++ {
		out.Chart = append(out.Chart, EMIPoint{
			Year:           y,
			Balance:        balanceAtYear(base, y),
			PrepaidBalance: balanceAtYear(prepaid, y),
		})
	}
	return out
}

func simulateLoan(principal, r, payment float64, maxMonths int) loanSim {
	sim := loanSim{opening: principal, monthly: make([]float64, 0, maxMonths)}
	if principal <= 0.01 {
		sim.months = 0
		return sim
	}
	if maxMonths < 1 {
		maxMonths = 1
	}
	bal := principal
	const eps = 0.01
	for m := 1; m <= maxMonths; m++ {
		intM := bal * r
		due := bal + intM
		pay := payment
		if pay <= 0 {
			sim.interest += intM
			bal = due
			sim.monthly = append(sim.monthly, bal)
			sim.months = m
			continue
		}
		if pay > due {
			pay = due
		}
		sim.interest += intM
		bal = due - pay
		if bal <= eps {
			bal = 0
			sim.monthly = append(sim.monthly, 0)
			sim.months = m
			return sim
		}
		sim.monthly = append(sim.monthly, bal)
		sim.months = m
	}
	return sim
}

func balanceAtYear(sim loanSim, year int) float64 {
	if year <= 0 {
		return sim.opening
	}
	idx := year*12 - 1
	if idx >= len(sim.monthly) {
		if len(sim.monthly) == 0 {
			return 0
		}
		last := sim.monthly[len(sim.monthly)-1]
		if last <= 0.01 {
			return 0
		}
		return last
	}
	return sim.monthly[idx]
}
