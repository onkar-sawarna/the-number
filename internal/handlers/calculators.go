package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thenumber/app/internal/calc"
	"github.com/thenumber/app/web/templates"
)

func (s *Server) firePage(c echo.Context) error {
	return render(c, http.StatusOK, templates.FIREPage(s.page(c, "FIRE calculator")))
}

func (s *Server) sipPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.SIPPage(s.page(c, "SIP calculator")))
}

func (s *Server) emiPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.EMIPage(s.page(c, "EMI calculator")))
}

func (s *Server) emergencyPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.EmergencyPage(s.page(c, "Emergency fund")))
}

func (s *Server) budgetPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.BudgetPage(s.page(c, "50/30/20 budget")))
}

func (s *Server) saveFIRE(c echo.Context) error {
	d := calc.DefaultFIRE()
	in := calc.FIREInput{
		Age:            parseInt(c, "age", d.Age),
		AnnualExpenses: parseFloat(c, "annual_expenses", d.AnnualExpenses),
		CurrentSavings: parseFloat(c, "current_savings", d.CurrentSavings),
		MonthlySavings: parseFloat(c, "monthly_savings", d.MonthlySavings),
		ExpectedReturn: parseFloat(c, "expected_return", d.ExpectedReturn),
		Inflation:      parseFloat(c, "inflation", d.Inflation),
		SWR:            parseFloat(c, "swr", d.SWR),
	}
	out := calc.FIRE(in)
	return s.saveScenario(c, "fire", titleOr(c, fmt.Sprintf("FIRE · age %d", in.Age)), in, out)
}

func (s *Server) saveSIP(c echo.Context) error {
	d := calc.DefaultSIP()
	in := calc.SIPInput{
		Monthly:        parseFloat(c, "monthly", d.Monthly),
		Existing:       parseFloat(c, "existing", d.Existing),
		ExpectedReturn: parseFloat(c, "expected_return", d.ExpectedReturn),
		Years:          parseInt(c, "years", d.Years),
	}
	out := calc.SIP(in)
	return s.saveScenario(c, "sip", titleOr(c, fmt.Sprintf("SIP · %d years", in.Years)), in, out)
}

func (s *Server) saveEMI(c echo.Context) error {
	d := calc.DefaultEMI()
	in := calc.EMIInput{
		Principal:    parseFloat(c, "principal", d.Principal),
		AnnualRate:   parseFloat(c, "annual_rate", d.AnnualRate),
		Years:        parseInt(c, "years", d.Years),
		ExtraMonthly: parseFloat(c, "extra_monthly", d.ExtraMonthly),
		Lump:         parseFloat(c, "lump", d.Lump),
	}
	out := calc.EMI(in)
	return s.saveScenario(c, "emi", titleOr(c, fmt.Sprintf("EMI · %s", calc.FormatINR(in.Principal))), in, out)
}

func (s *Server) saveEmergency(c echo.Context) error {
	d := calc.DefaultEmergency()
	in := calc.EmergencyInput{
		MonthlyEssentials: parseFloat(c, "monthly_essentials", d.MonthlyEssentials),
		MonthsCover:       parseFloat(c, "months_cover", d.MonthsCover),
		CurrentBuffer:     parseFloat(c, "current_buffer", d.CurrentBuffer),
		MonthlyTopup:      parseFloat(c, "monthly_topup", d.MonthlyTopup),
		ParkingReturn:     parseFloat(c, "parking_return", d.ParkingReturn),
	}
	out := calc.Emergency(in)
	return s.saveScenario(c, "emergency", titleOr(c, fmt.Sprintf("Emergency · %.0f months", in.MonthsCover)), in, out)
}

func (s *Server) saveBudget(c echo.Context) error {
	d := calc.DefaultBudget()
	in := calc.BudgetInput{
		Income:  parseFloat(c, "income", d.Income),
		Needs:   parseFloat(c, "needs", d.Needs),
		Wants:   parseFloat(c, "wants", d.Wants),
		Savings: parseFloat(c, "savings", d.Savings),
	}
	out := calc.Budget(in)
	return s.saveScenario(c, "budget", titleOr(c, "50/30/20 budget"), in, out)
}
