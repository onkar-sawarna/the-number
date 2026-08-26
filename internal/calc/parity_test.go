package calc

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

type jsJob struct {
	Op    string  `json:"op"`
	Input any     `json:"input,omitempty"`
	N     float64 `json:"n,omitempty"`
}

func TestGoMatchesCalcJS(t *testing.T) {
	jsResults := runCalcJS(t, parityJobs())

	i := 0
	next := func() json.RawMessage {
		if i >= len(jsResults) {
			t.Fatal("ran out of JS results")
		}
		v := jsResults[i]
		i++
		return v
	}

	for _, in := range fireFixtures() {
		var js jsFIRE
		mustUnmarshal(t, next(), &js)
		goOut := FIRE(in)
		compareFIRE(t, in, js, goOut)
	}

	for _, in := range sipFixtures() {
		var js jsSIP
		mustUnmarshal(t, next(), &js)
		goOut := SIP(in)
		if math.Abs(js.Invested-goOut.Invested) > 0.02 || math.Abs(js.FV-goOut.FV) > 0.02 || math.Abs(js.Gain-goOut.Gain) > 0.02 {
			t.Errorf("SIP %+v: js invested/fv/gain=%v/%v/%v go=%v/%v/%v", in, js.Invested, js.FV, js.Gain, goOut.Invested, goOut.FV, goOut.Gain)
		}
		if len(js.Chart) != len(goOut.Chart) {
			t.Errorf("SIP chart len js=%d go=%d", len(js.Chart), len(goOut.Chart))
		} else if len(js.Chart) > 0 {
			last := len(js.Chart) - 1
			if math.Abs(js.Chart[last].FV-goOut.Chart[last].FV) > 0.02 {
				t.Errorf("SIP last FV js=%v go=%v", js.Chart[last].FV, goOut.Chart[last].FV)
			}
		}
	}

	for _, in := range emiFixtures() {
		var js jsEMI
		mustUnmarshal(t, next(), &js)
		goOut := EMI(in)
		if math.Abs(js.EMI-goOut.EMI) > 0.02 || js.Months != goOut.Months || js.MonthsSaved != goOut.MonthsSaved {
			t.Errorf("EMI %+v: js emi/months/saved=%v/%d/%d go=%v/%d/%d", in, js.EMI, js.Months, js.MonthsSaved, goOut.EMI, goOut.Months, goOut.MonthsSaved)
		}
		if math.Abs(js.TotalInterest-goOut.TotalInterest) > 0.02 || math.Abs(js.InterestSaved-goOut.InterestSaved) > 0.02 {
			t.Errorf("EMI interest js=%v saved=%v go=%v saved=%v", js.TotalInterest, js.InterestSaved, goOut.TotalInterest, goOut.InterestSaved)
		}
		if len(js.Chart) != len(goOut.Chart) {
			t.Errorf("EMI chart len js=%d go=%d", len(js.Chart), len(goOut.Chart))
		}
	}

	for _, in := range emergencyFixtures() {
		var js jsEmergency
		mustUnmarshal(t, next(), &js)
		goOut := Emergency(in)
		if math.Abs(js.Target-goOut.Target) > 0.02 || math.Abs(js.Gap-goOut.Gap) > 0.02 || js.Reaches != goOut.Reaches {
			t.Errorf("emergency %+v: js target/gap/reaches=%v/%v/%v go=%v/%v/%v", in, js.Target, js.Gap, js.Reaches, goOut.Target, goOut.Gap, goOut.Reaches)
		}
		if math.Abs(js.MonthsToFill-goOut.MonthsToFill) > 0.02 {
			t.Errorf("emergency months js=%v go=%v", js.MonthsToFill, goOut.MonthsToFill)
		}
		if len(js.Chart) != len(goOut.Chart) {
			t.Errorf("emergency chart len js=%d go=%d", len(js.Chart), len(goOut.Chart))
		}
	}

	for _, in := range budgetFixtures() {
		var js jsBudget
		mustUnmarshal(t, next(), &js)
		goOut := Budget(in)
		if math.Abs(js.TargetNeeds-goOut.TargetNeeds) > 0.02 || math.Abs(js.Unallocated-goOut.Unallocated) > 0.02 || js.Overspent != goOut.Overspent {
			t.Errorf("budget %+v js=%+v go=%+v", in, js, goOut)
		}
		if math.Abs(js.SavingsRate-goOut.SavingsRate) > 1e-9 {
			t.Errorf("budget savingsRate js=%v go=%v", js.SavingsRate, goOut.SavingsRate)
		}
	}

	for _, n := range formatINRFixtures() {
		var js string
		mustUnmarshal(t, next(), &js)
		if got := FormatINR(n); js != got {
			t.Errorf("formatINR(%v): js=%q go=%q", n, js, got)
		}
	}
	for _, n := range formatUSDFixtures() {
		var js string
		mustUnmarshal(t, next(), &js)
		if got := FormatUSD(n); js != got {
			t.Errorf("formatUSD(%v): js=%q go=%q", n, js, got)
		}
	}

	if i != len(jsResults) {
		t.Fatalf("unused JS results: consumed %d of %d", i, len(jsResults))
	}
}

func parityJobs() []jsJob {
	var jobs []jsJob
	for _, in := range fireFixtures() {
		jobs = append(jobs, jsJob{Op: "fire", Input: fireToJS(in)})
	}
	for _, in := range sipFixtures() {
		jobs = append(jobs, jsJob{Op: "sip", Input: map[string]any{
			"monthly": in.Monthly, "existing": in.Existing, "expectedReturn": in.ExpectedReturn, "years": in.Years,
		}})
	}
	for _, in := range emiFixtures() {
		jobs = append(jobs, jsJob{Op: "emi", Input: map[string]any{
			"principal": in.Principal, "annualRate": in.AnnualRate, "years": in.Years, "extraMonthly": in.ExtraMonthly, "lump": in.Lump,
		}})
	}
	for _, in := range emergencyFixtures() {
		jobs = append(jobs, jsJob{Op: "emergency", Input: map[string]any{
			"monthlyEssentials": in.MonthlyEssentials, "monthsCover": in.MonthsCover, "currentBuffer": in.CurrentBuffer,
			"monthlyTopup": in.MonthlyTopup, "parkingReturn": in.ParkingReturn,
		}})
	}
	for _, in := range budgetFixtures() {
		jobs = append(jobs, jsJob{Op: "budget", Input: map[string]any{
			"income": in.Income, "needs": in.Needs, "wants": in.Wants, "savings": in.Savings,
		}})
	}
	for _, n := range formatINRFixtures() {
		jobs = append(jobs, jsJob{Op: "formatINR", N: n})
	}
	for _, n := range formatUSDFixtures() {
		jobs = append(jobs, jsJob{Op: "formatUSD", N: n})
	}
	return jobs
}

func fireFixtures() []FIREInput {
	buy := DefaultFIRE()
	buy.Housing = "buy"
	buy.CityTier = 1

	us := FIREInput{
		Age:            32,
		AnnualExpenses: 60_000,
		CurrentSavings: 80_000,
		MonthlySavings: 2_500,
		ExpectedReturn: 8,
		Inflation:      3,
		SWR:            4,
		NPSNow:         40_000,
		NPSMonthly:     500,
		StepUp:         5,
		CityTier:       2,
		Housing:        "rent",
		Region:         "us",
	}

	pots := DefaultFIRE()
	pots.NPSNow = 1_000_000
	pots.NPSMonthly = 10_000
	pots.PPFNow = 200_000
	pots.PPFMonthly = 12_500
	pots.EPFNow = 500_000
	pots.EPFMonthly = 8_000
	pots.ForeignNow = 100_000
	pots.ForeignMonthly = 5_000
	pots.StoppedNow = 200_000
	pots.GoldNow = 150_000
	pots.GoldMonthly = 2_000
	pots.JewelleryNow = 300_000

	already := FIREInput{Age: 40, AnnualExpenses: 1_200_000, CurrentSavings: 40_000_000, ExpectedReturn: 8, Inflation: 6, SWR: 4, Housing: "rent"}
	never := FIREInput{Age: 30, AnnualExpenses: 1_200_000, ExpectedReturn: 0, Inflation: 0, SWR: 4, Housing: "rent"}
	mixed := DefaultFIRE()
	mixed.Mixed = true
	mixed.UsdParked = 10_000
	mixed.UsdMonthly = 500
	return []FIREInput{DefaultFIRE(), buy, us, pots, already, never, mixed}
}

func sipFixtures() []SIPInput {
	return []SIPInput{
		DefaultSIP(),
		{Monthly: 10_000, Existing: 0, ExpectedReturn: 12, Years: 10},
		{Monthly: 10_000, Existing: 50_000, ExpectedReturn: 0, Years: 10},
	}
}

func emiFixtures() []EMIInput {
	return []EMIInput{
		DefaultEMI(),
		{Principal: 1_000_000, AnnualRate: 10, Years: 20, ExtraMonthly: 0, Lump: 0},
		{Principal: 1_000_000, AnnualRate: 10, Years: 20, ExtraMonthly: 2_000, Lump: 50_000},
	}
}

func emergencyFixtures() []EmergencyInput {
	return []EmergencyInput{
		DefaultEmergency(),
		{MonthlyEssentials: 40_000, MonthsCover: 6, CurrentBuffer: 0, MonthlyTopup: 20_000, ParkingReturn: 0},
	}
}

func budgetFixtures() []BudgetInput {
	return []BudgetInput{
		DefaultBudget(),
		{Income: 100_000, Needs: 50_000, Wants: 30_000, Savings: 20_000},
		{Income: 100_000, Needs: 70_000, Wants: 40_000, Savings: 20_000},
	}
}

func formatINRFixtures() []float64 {
	return []float64{3e7, 1_250_000, 5_000, 500, -12_500}
}

func formatUSDFixtures() []float64 {
	return []float64{1_200_000, 45_000, 5_000, 500, -2_000}
}

func fireToJS(in FIREInput) map[string]any {
	return map[string]any{
		"age": in.Age, "annualExpenses": in.AnnualExpenses, "currentSavings": in.CurrentSavings,
		"monthlySavings": in.MonthlySavings, "expectedReturn": in.ExpectedReturn, "inflation": in.Inflation, "swr": in.SWR,
		"npsNow": in.NPSNow, "npsMonthly": in.NPSMonthly, "ppfNow": in.PPFNow, "ppfMonthly": in.PPFMonthly,
		"epfNow": in.EPFNow, "epfMonthly": in.EPFMonthly, "foreignNow": in.ForeignNow, "foreignMonthly": in.ForeignMonthly,
		"stoppedNow": in.StoppedNow, "goldNow": in.GoldNow, "goldMonthly": in.GoldMonthly, "jewelleryNow": in.JewelleryNow,
		"stepUp": in.StepUp, "cityTier": in.CityTier, "housing": in.Housing, "region": in.Region,
		"mixed": in.Mixed, "fxINRPerUSD": in.FxINRPerUSD,
		"usdParked": in.UsdParked, "usdMonthly": in.UsdMonthly, "usdStopped": in.UsdStopped, "usdRetire": in.UsdRetire,
		"indiaParked": in.IndiaParked, "indiaNpsNow": in.IndiaNPSNow, "indiaNpsMonthly": in.IndiaNPSMonthly,
		"indiaGold": in.IndiaGold, "indiaJew": in.IndiaJew,
		"contribScale": in.ContribScale, "pauseMonths": in.PauseMonths,
	}
}

type jsFIRE struct {
	FireNumber      float64 `json:"fireNumber"`
	FireNumberLater float64 `json:"fireNumberLater"`
	Lifestyle       float64 `json:"lifestyle"`
	LifestyleLater  float64 `json:"lifestyleLater"`
	HouseAdd        float64 `json:"houseAdd"`
	StartingCorpus  float64 `json:"startingCorpus"`
	Jewellery       float64 `json:"jewellery"`
	JewelleryLater  float64 `json:"jewelleryLater"`
	MonthlyIn       float64 `json:"monthlyIn"`
	Lean            float64 `json:"lean"`
	Regular         float64 `json:"regular"`
	Fat             float64 `json:"fat"`
	Years           float64 `json:"years"`
	ReachesFire     bool    `json:"reachesFire"`
	FIAge           int     `json:"fiAge"`
	CrossingYears   float64 `json:"crossingYears"`
	ReachesCrossing bool    `json:"reachesCrossing"`
	CrossingAge     int     `json:"crossingAge"`
	Chart           []struct {
		Year     int     `json:"year"`
		Corpus   float64 `json:"corpus"`
		NetWorth float64 `json:"netWorth"`
		Target   float64 `json:"target"`
	} `json:"chart"`
}

type jsSIP struct {
	Invested float64 `json:"invested"`
	FV       float64 `json:"fv"`
	Gain     float64 `json:"gain"`
	Chart    []struct {
		Year int     `json:"year"`
		FV   float64 `json:"fv"`
	} `json:"chart"`
}

type jsEMI struct {
	EMI           float64 `json:"emi"`
	Months        int     `json:"months"`
	TotalInterest float64 `json:"totalInterest"`
	InterestSaved float64 `json:"interestSaved"`
	MonthsSaved   int     `json:"monthsSaved"`
	Chart         []any   `json:"chart"`
}

type jsEmergency struct {
	Target       float64 `json:"target"`
	Gap          float64 `json:"gap"`
	MonthsToFill float64 `json:"monthsToFill"`
	Reaches      bool    `json:"reaches"`
	Chart        []any   `json:"chart"`
}

type jsBudget struct {
	TargetNeeds float64 `json:"targetNeeds"`
	Unallocated float64 `json:"unallocated"`
	Overspent   bool    `json:"overspent"`
	SavingsRate float64 `json:"savingsRate"`
}

func compareFIRE(t *testing.T, in FIREInput, js jsFIRE, goOut FIREOutput) {
	t.Helper()
	label := fireLabel(in)
	same := func(name string, a, b float64) {
		if math.Abs(a-b) > 0.02 && math.Abs(a-b)/math.Max(1, math.Abs(b)) > 1e-9 {
			t.Errorf("%s %s: js=%v go=%v", label, name, a, b)
		}
	}
	same("fireNumber", js.FireNumber, goOut.FireNumber)
	same("fireNumberLater", js.FireNumberLater, goOut.FireNumberLater)
	same("lifestyle", js.Lifestyle, goOut.Lifestyle)
	same("houseAdd", js.HouseAdd, goOut.HouseAdd)
	same("startingCorpus", js.StartingCorpus, goOut.StartingCorpus)
	same("jewelleryLater", js.JewelleryLater, goOut.JewelleryLater)
	same("monthlyIn", js.MonthlyIn, goOut.MonthlyIn)
	same("lean", js.Lean, goOut.Lean)
	same("fat", js.Fat, goOut.Fat)
	same("years", js.Years, goOut.Years)
	same("crossingYears", js.CrossingYears, goOut.CrossingYears)
	if js.ReachesFire != goOut.ReachesFire || js.FIAge != goOut.FIAge {
		t.Errorf("%s reaches/fiAge js=%v/%d go=%v/%d", label, js.ReachesFire, js.FIAge, goOut.ReachesFire, goOut.FIAge)
	}
	if js.ReachesCrossing != goOut.ReachesCrossing || js.CrossingAge != goOut.CrossingAge {
		t.Errorf("%s crossing js=%v/%d go=%v/%d", label, js.ReachesCrossing, js.CrossingAge, goOut.ReachesCrossing, goOut.CrossingAge)
	}
	if len(js.Chart) != len(goOut.Chart) {
		t.Errorf("%s chart len js=%d go=%d", label, len(js.Chart), len(goOut.Chart))
		return
	}
	if len(js.Chart) == 0 {
		return
	}
	same("chart[0].corpus", js.Chart[0].Corpus, goOut.Chart[0].Corpus)
	last := len(js.Chart) - 1
	same("chart[last].corpus", js.Chart[last].Corpus, goOut.Chart[last].Corpus)
	same("chart[last].target", js.Chart[last].Target, goOut.Chart[last].Target)
	same("chart[last].netWorth", js.Chart[last].NetWorth, goOut.Chart[last].NetWorth)
}

func fireLabel(in FIREInput) string {
	b, _ := json.Marshal(map[string]any{"region": in.Region, "housing": in.Housing, "expenses": in.AnnualExpenses, "savings": in.CurrentSavings})
	return string(b)
}

func runCalcJS(t *testing.T, jobs []jsJob) []json.RawMessage {
	t.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	root := filepath.Join(dir, "../..")
	bridge := filepath.Join(dir, "testdata/calc_bridge.js")
	calcJS := filepath.Join(root, "web/static/js/calc.js")

	payload, err := json.Marshal(jobs)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", bridge, calcJS)
	cmd.Stdin = bytes.NewReader(payload)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("node bridge: %v\n%s", err, stderr.String())
	}
	var results []json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		t.Fatalf("decode JS results: %v\n%s", err, stdout.String())
	}
	return results
}

func mustUnmarshal(t *testing.T, raw json.RawMessage, dest any) {
	t.Helper()
	if err := json.Unmarshal(raw, dest); err != nil {
		t.Fatalf("unmarshal JS result: %v (%s)", err, raw)
	}
}
