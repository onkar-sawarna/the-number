package calc

import "math"

// SIPBump is the extra monthly SIP used for “what moves the year” (primary currency).
func SIPBump(region string) float64 {
	if IsUSD(region) {
		return 100
	}
	return 5_000
}

// YearShift is one FIRE run against the current plan.
type YearShift struct {
	Kind       string  `json:"kind"`
	Amount     float64 `json:"amount"`
	Years      float64 `json:"years"`
	Reaches    bool    `json:"reaches"`
	FIAge      int     `json:"fi_age"`
	DeltaYears float64 `json:"delta_years"`
}

func shiftFrom(base, out FIREOutput, kind string, amount float64) YearShift {
	s := YearShift{
		Kind:    kind,
		Amount:  amount,
		Years:   out.Years,
		Reaches: out.ReachesFire,
		FIAge:   out.FIAge,
	}
	if base.ReachesFire && out.ReachesFire {
		s.DeltaYears = out.Years - base.Years
	} else if base.ReachesFire && !out.ReachesFire {
		s.DeltaYears = 80
	} else if !base.ReachesFire && out.ReachesFire {
		s.DeltaYears = out.Years - 80
	}
	return s
}

// YearMoves: extra SIP, 1% lower return, buy vs rent the house.
func YearMoves(in FIREInput) []YearShift {
	base := FIRE(in)
	sip := in
	sip.MonthlySavings += SIPBump(in.Region)
	ret := in
	ret.ExpectedReturn = math.Max(0, in.ExpectedReturn-1)
	house := in
	if in.Housing == "buy" {
		house.Housing = "rent"
	} else {
		house.Housing = "buy"
	}
	return []YearShift{
		shiftFrom(base, FIRE(sip), "sip", SIPBump(in.Region)),
		shiftFrom(base, FIRE(ret), "return", 1),
		shiftFrom(base, FIRE(house), "house", 0),
	}
}

// YearOptions: 40% less going in, half going in, two years off.
func YearOptions(in FIREInput) []YearShift {
	base := FIRE(in)
	cut := in
	cut.ContribScale = 0.6
	part := in
	part.ContribScale = 0.5
	pause := in
	pause.PauseMonths = 24
	return []YearShift{
		shiftFrom(base, FIRE(cut), "paycut", 40),
		shiftFrom(base, FIRE(part), "parttime", 50),
		shiftFrom(base, FIRE(pause), "pause", 24),
	}
}
