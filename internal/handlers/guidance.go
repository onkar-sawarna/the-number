package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thenumber/app/internal/ai"
	"github.com/thenumber/app/web/templates"
)

func (s *Server) guidanceGet(c echo.Context) error {
	return render(c, http.StatusOK, templates.GuidancePage(s.page(c, "FIRE sleeves", "title_guidance")))
}

func (s *Server) guidancePost(c echo.Context) error {
	risk := strings.ToLower(strings.TrimSpace(c.FormValue("risk")))
	if risk != "conservative" && risk != "aggressive" {
		risk = "moderate"
	}
	in := ai.GuideInput{
		Age:               parseInt(c, "age", 30),
		Horizon:           parseInt(c, "horizon", 10),
		Risk:              risk,
		MonthlyInvestable: parseFloat(c, "monthly_investable", 25000),
		Goals:             strings.TrimSpace(c.FormValue("goals")),
	}
	res := ai.Guide(c.Request().Context(), in)

	labels := make([]string, 0, len(res.Sleeves))
	values := make([]string, 0, len(res.Sleeves))
	sleeves := make([]templates.SleeveVM, 0, len(res.Sleeves))
	for _, sl := range res.Sleeves {
		if sl.Percent <= 0 {
			continue
		}
		labels = append(labels, sl.Category)
		values = append(values, fmt.Sprintf("%.1f", sl.Percent))
		sleeves = append(sleeves, templates.SleeveVM{
			Category: sl.Category,
			Percent:  strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%.1f", sl.Percent), "0"), "."),
		})
	}
	return render(c, http.StatusOK, templates.GuidanceResult(templates.GuidanceResultView{
		Warning:   res.Warning,
		Rationale: res.Rationale,
		Source:    res.Source,
		Labels:    strings.Join(labels, "|"),
		Values:    strings.Join(values, "|"),
		Sleeves:   sleeves,
	}))
}
