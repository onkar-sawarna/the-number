package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thenumber/app/internal/calc"
	"github.com/thenumber/app/internal/models"
	"github.com/thenumber/app/web/templates"
)

func (s *Server) dashboard(c echo.Context) error {
	return render(c, http.StatusOK, templates.DashboardPage(s.page(c, "Dashboard")))
}

func (s *Server) listRows(userID uint) ([]templates.ScenarioVM, error) {
	var rows []models.Scenario
	if err := s.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]templates.ScenarioVM, 0, len(rows))
	for _, r := range rows {
		out = append(out, templates.ScenarioVM{
			ID:         r.ID,
			Date:       r.CreatedAt.Local().Format("2 Jan 2006"),
			Kind:       r.Kind,
			KindLabel:  kindLabel(r.Kind),
			Title:      r.Title,
			DeletePath: fmt.Sprintf("/scenarios/%d", r.ID),
		})
	}
	return out, nil
}

func (s *Server) dashboardList(c echo.Context) error {
	u := s.currentUser(c)
	rows, err := s.listRows(u.ID)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `<p>Could not load scenarios.</p>`)
	}
	return render(c, http.StatusOK, templates.ScenarioList(rows))
}

func (s *Server) deleteScenario(c echo.Context) error {
	u := s.currentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res := s.DB.Where("id = ? AND user_id = ?", id, u.ID).Delete(&models.Scenario{})
	if res.Error != nil {
		return c.HTML(http.StatusInternalServerError, `<p>Could not delete.</p>`)
	}
	rows, err := s.listRows(u.ID)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `<p>Could not load scenarios.</p>`)
	}
	return render(c, http.StatusOK, templates.ScenarioList(rows))
}

func (s *Server) dashboardCompare(c echo.Context) error {
	p := s.page(c, "Compare")
	ids := c.QueryParams()["id"]
	if len(ids) != 2 {
		return render(c, http.StatusBadRequest, templates.CompareError(p, "Compare needs exactly two scenarios."))
	}
	u := s.currentUser(c)
	cards := make([]templates.CompareCard, 0, 2)
	for _, raw := range ids {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return render(c, http.StatusBadRequest, templates.CompareError(p, "Compare needs exactly two scenarios."))
		}
		var sc models.Scenario
		if err := s.DB.Where("id = ? AND user_id = ?", id, u.ID).First(&sc).Error; err != nil {
			return render(c, http.StatusNotFound, templates.CompareError(p, "One of those scenarios was not found."))
		}
		cards = append(cards, compareCard(sc))
	}
	return render(c, http.StatusOK, templates.ComparePage(p, cards))
}

func compareCard(sc models.Scenario) templates.CompareCard {
	card := templates.CompareCard{
		Title: sc.Title,
		Kind:  kindLabel(sc.Kind),
		Date:  sc.CreatedAt.Local().Format(time.DateOnly),
	}
	switch sc.Kind {
	case "fire":
		var o calc.FIREOutput
		_ = json.Unmarshal([]byte(sc.Outputs), &o)
		years := "Does not reach FIRE within 80 years"
		if o.ReachesFire && o.Years == 0 {
			years = "Already independent"
		} else if o.ReachesFire {
			years = fmt.Sprintf("%.1f years", o.Years)
		}
		card.Rows = []templates.KV{
			{Label: "FIRE number", Value: calc.FormatINR(o.FireNumber)},
			{Label: "Years", Value: years},
			{Label: "Lean", Value: calc.FormatINR(o.Lean)},
			{Label: "Regular", Value: calc.FormatINR(o.Regular)},
			{Label: "Fat", Value: calc.FormatINR(o.Fat)},
			{Label: "FI age", Value: fmt.Sprintf("%d", o.FIAge)},
		}
	case "sip":
		var o calc.SIPOutput
		_ = json.Unmarshal([]byte(sc.Outputs), &o)
		card.Rows = []templates.KV{
			{Label: "Invested", Value: calc.FormatINR(o.Invested)},
			{Label: "Future value", Value: calc.FormatINR(o.FV)},
			{Label: "Gain", Value: calc.FormatINR(o.Gain)},
		}
	case "emi":
		var o calc.EMIOutput
		_ = json.Unmarshal([]byte(sc.Outputs), &o)
		card.Rows = []templates.KV{
			{Label: "EMI", Value: calc.FormatINR(o.EMI)},
			{Label: "Months", Value: fmt.Sprintf("%d", o.Months)},
			{Label: "Total interest", Value: calc.FormatINR(o.TotalInterest)},
			{Label: "Interest saved", Value: calc.FormatINR(o.InterestSaved)},
			{Label: "Months saved", Value: fmt.Sprintf("%d", o.MonthsSaved)},
		}
	case "emergency":
		var o calc.EmergencyOutput
		_ = json.Unmarshal([]byte(sc.Outputs), &o)
		fill := "Does not fill within 40 years"
		if o.Reaches && o.Gap == 0 {
			fill = "Already funded"
		} else if o.Reaches {
			fill = fmt.Sprintf("%.0f months", o.MonthsToFill)
		}
		card.Rows = []templates.KV{
			{Label: "Target", Value: calc.FormatINR(o.Target)},
			{Label: "Gap", Value: calc.FormatINR(o.Gap)},
			{Label: "Coverage now", Value: fmt.Sprintf("%.1f months", o.CoverageNow)},
			{Label: "Months to fill", Value: fill},
		}
	case "budget":
		var o calc.BudgetOutput
		_ = json.Unmarshal([]byte(sc.Outputs), &o)
		flag := "No"
		if o.Overspent {
			flag = "Yes"
		}
		card.Rows = []templates.KV{
			{Label: "Needs target", Value: calc.FormatINR(o.TargetNeeds)},
			{Label: "Wants target", Value: calc.FormatINR(o.TargetWants)},
			{Label: "Savings target", Value: calc.FormatINR(o.TargetSavings)},
			{Label: "Unallocated", Value: calc.FormatINR(o.Unallocated)},
			{Label: "Savings rate", Value: fmt.Sprintf("%.1f%%", o.SavingsRate)},
			{Label: "Overspent", Value: flag},
		}
	}
	return card
}
