package calc

type EmergencyInput struct {
	MonthlyEssentials float64 `json:"monthly_essentials"`
	MonthsCover       float64 `json:"months_cover"`
	CurrentBuffer     float64 `json:"current_buffer"`
	MonthlyTopup      float64 `json:"monthly_topup"`
	ParkingReturn     float64 `json:"parking_return"`
}

type EmergencyPoint struct {
	Month   int     `json:"month"`
	Balance float64 `json:"balance"`
	Target  float64 `json:"target"`
}

type EmergencyOutput struct {
	Target       float64          `json:"target"`
	Gap          float64          `json:"gap"`
	CoverageNow  float64          `json:"coverage_now"`
	MonthsToFill float64          `json:"months_to_fill"`
	Reaches      bool             `json:"reaches"`
	Chart        []EmergencyPoint `json:"chart"`
}

func DefaultEmergency() EmergencyInput {
	return EmergencyInput{
		MonthlyEssentials: 60_000,
		MonthsCover:       6,
		CurrentBuffer:     80_000,
		MonthlyTopup:      15_000,
		ParkingReturn:     6,
	}
}

func Emergency(in EmergencyInput) EmergencyOutput {
	target := in.MonthlyEssentials * in.MonthsCover
	gap := target - in.CurrentBuffer
	if gap < 0 {
		gap = 0
	}
	cov := 0.0
	if in.MonthlyEssentials > 0 {
		cov = in.CurrentBuffer / in.MonthlyEssentials
	}
	out := EmergencyOutput{Target: target, Gap: gap, CoverageNow: cov}

	if gap == 0 {
		out.MonthsToFill = 0
		out.Reaches = true
		out.Chart = emergencyChart(in, target, 0)
		return out
	}

	rm := in.ParkingReturn / 12.0 / 100.0
	bal := in.CurrentBuffer
	const maxM = 40 * 12
	for m := 1; m <= maxM; m++ {
		bal = bal*(1+rm) + in.MonthlyTopup
		if bal >= target {
			out.MonthsToFill = float64(m)
			out.Reaches = true
			out.Chart = emergencyChart(in, target, m)
			return out
		}
	}
	out.Reaches = false
	out.MonthsToFill = 0
	out.Chart = emergencyChart(in, target, maxM)
	return out
}

func emergencyChart(in EmergencyInput, target float64, fillMonths int) []EmergencyPoint {
	horizon := fillMonths
	if horizon <= 0 {
		horizon = 12
	}
	if horizon < 12 {
		horizon = 12
	}
	// yearly-ish snapshots plus start; cap at 40 years of months.
	if horizon > 40*12 {
		horizon = 40 * 12
	}
	rm := in.ParkingReturn / 12.0 / 100.0
	bal := in.CurrentBuffer
	pts := []EmergencyPoint{{Month: 0, Balance: bal, Target: target}}
	step := 1
	if horizon > 36 {
		step = 12
	}
	for m := 1; m <= horizon; m++ {
		bal = bal*(1+rm) + in.MonthlyTopup
		if m%step == 0 || m == horizon || (fillMonths > 0 && m == fillMonths) {
			pts = append(pts, EmergencyPoint{Month: m, Balance: bal, Target: target})
		}
	}
	return pts
}
