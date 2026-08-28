package calc

type SIPInput struct {
	Monthly        float64 `json:"monthly"`
	Existing       float64 `json:"existing"`
	ExpectedReturn float64 `json:"expected_return"`
	Years          int     `json:"years"`
	StepUp         float64 `json:"step_up"` // yearly % rise in monthly SIP
}

type SIPPoint struct {
	Year     int     `json:"year"`
	Invested float64 `json:"invested"`
	FV       float64 `json:"fv"`
}

type SIPOutput struct {
	Invested float64    `json:"invested"`
	FV       float64    `json:"fv"`
	Gain     float64    `json:"gain"`
	Chart    []SIPPoint `json:"chart"`
}

func DefaultSIP() SIPInput {
	return SIPInput{
		Monthly:        10_000,
		Existing:       0,
		ExpectedReturn: 12,
		Years:          15,
		StepUp:         10,
	}
}

func SIP(in SIPInput) SIPOutput {
	if in.Years < 0 {
		in.Years = 0
	}
	n := in.Years * 12
	fv, invested := sipPath(in.Monthly, in.Existing, in.StepUp, in.ExpectedReturn, n)
	out := SIPOutput{
		Invested: invested,
		FV:       fv,
		Gain:     fv - invested,
		Chart:    make([]SIPPoint, 0, in.Years+1),
	}
	for y := 0; y <= in.Years; y++ {
		fy, inv := sipPath(in.Monthly, in.Existing, in.StepUp, in.ExpectedReturn, y*12)
		out.Chart = append(out.Chart, SIPPoint{
			Year:     y,
			Invested: inv,
			FV:       fy,
		})
	}
	return out
}

func sipPath(monthly, existing, stepUp, annualPct float64, months int) (fv, invested float64) {
	if months < 0 {
		months = 0
	}
	rm := monthlyEffective(annualPct)
	fv = nz(existing)
	invested = nz(existing)
	for m := 1; m <= months; m++ {
		p := steppedMonthly(nz(monthly), stepUp, m)
		fv = fv*(1+rm) + p
		invested += p
	}
	return fv, invested
}
