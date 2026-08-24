package calc

import "math"

type SIPInput struct {
	Monthly        float64 `json:"monthly"`
	Existing       float64 `json:"existing"`
	ExpectedReturn float64 `json:"expected_return"`
	Years          int     `json:"years"`
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
	}
}

func SIP(in SIPInput) SIPOutput {
	if in.Years < 0 {
		in.Years = 0
	}
	n := in.Years * 12
	rm := monthlyEffective(in.ExpectedReturn)
	fv := sipFV(in.Monthly, in.Existing, rm, n)
	invested := in.Monthly*float64(n) + in.Existing
	out := SIPOutput{
		Invested: invested,
		FV:       fv,
		Gain:     fv - invested,
		Chart:    make([]SIPPoint, 0, in.Years+1),
	}
	for y := 0; y <= in.Years; y++ {
		nm := y * 12
		out.Chart = append(out.Chart, SIPPoint{
			Year:     y,
			Invested: in.Monthly*float64(nm) + in.Existing,
			FV:       sipFV(in.Monthly, in.Existing, rm, nm),
		})
	}
	return out
}

func sipFV(monthly, existing, rm float64, n int) float64 {
	if n <= 0 {
		return existing
	}
	nf := float64(n)
	if nearlyZero(rm) {
		return existing + monthly*nf
	}
	fvSip := monthly * (math.Pow(1+rm, nf) - 1) / rm
	fvExist := existing * math.Pow(1+rm, nf)
	return fvSip + fvExist
}
