package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
