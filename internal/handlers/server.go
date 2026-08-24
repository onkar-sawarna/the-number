package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/thenumber/app/internal/models"
	"github.com/thenumber/app/internal/session"
	"github.com/thenumber/app/web/templates"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Server {
	return &Server{DB: db}
}

func (s *Server) Register(e *echo.Echo) {
	e.GET("/", s.landing)
	e.GET("/disclaimer", s.disclaimer)
	e.POST("/theme", s.theme)

	e.GET("/calculators/fire", s.firePage)
	e.GET("/calculators/sip", s.sipPage)
	e.GET("/calculators/emi", s.emiPage)
	e.GET("/calculators/emergency", s.emergencyPage)
	e.GET("/calculators/budget", s.budgetPage)

	e.GET("/guidance", s.guidanceGet)
	e.POST("/guidance", s.guidancePost)
}

func render(c echo.Context, status int, t templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(status)
	return t.Render(c.Request().Context(), c.Response())
}

func (s *Server) page(c echo.Context, title string) templates.Page {
	return templates.Page{
		Title: title,
		Dark:  themeDark(c),
		Path:  c.Request().URL.Path,
	}
}

func themeDark(c echo.Context) bool {
	ck, err := c.Cookie("theme")
	return err == nil && ck.Value == "dark"
}

func (s *Server) currentUser(c echo.Context) *models.User {
	id, ok := session.UserID(c)
	if !ok {
		return nil
	}
	var u models.User
	if err := s.DB.First(&u, id).Error; err != nil {
		return nil
	}
	return &u
}

func (s *Server) requireUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.currentUser(c) == nil {
			if c.Request().Header.Get("HX-Request") == "true" {
				return c.HTML(http.StatusUnauthorized, `<p>Please log in to save.</p>`)
			}
			next := session.SafeNext(c.Request().URL.RequestURI())
			if next == "" {
				next = "/dashboard"
			}
			return c.Redirect(http.StatusSeeOther, "/login?next="+url.QueryEscape(next))
		}
		return next(c)
	}
}

func (s *Server) landing(c echo.Context) error {
	return render(c, http.StatusOK, templates.Landing(s.page(c, "Home")))
}

func (s *Server) disclaimer(c echo.Context) error {
	return render(c, http.StatusOK, templates.Disclaimer(s.page(c, "Disclaimer")))
}

func (s *Server) theme(c echo.Context) error {
	theme := c.FormValue("theme")
	if theme != "dark" && theme != "light" {
		theme = "light"
	}
	c.SetCookie(&http.Cookie{
		Name:     "theme",
		Value:    theme,
		Path:     "/",
		MaxAge:   31536000,
		SameSite: http.SameSiteLaxMode,
	})
	return c.NoContent(http.StatusNoContent)
}

func parseFloat(c echo.Context, name string, def float64) float64 {
	s := strings.TrimSpace(strings.ReplaceAll(c.FormValue(name), ",", ""))
	if s == "" {
		return def
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return v
}

func parseInt(c echo.Context, name string, def int) int {
	s := strings.TrimSpace(c.FormValue(name))
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func titleOr(c echo.Context, fallback string) string {
	t := strings.TrimSpace(c.FormValue("title"))
	if t == "" {
		return fallback
	}
	return t
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func (s *Server) saveScenario(c echo.Context, kind, title string, in any, out any) error {
	u := s.currentUser(c)
	if u == nil {
		return render(c, http.StatusUnauthorized, templates.LoginToSave())
	}
	sc := models.Scenario{
		UserID:  u.ID,
		Kind:    kind,
		Title:   title,
		Inputs:  mustJSON(in),
		Outputs: mustJSON(out),
	}
	if err := s.DB.Create(&sc).Error; err != nil {
		return c.HTML(http.StatusInternalServerError, `<p>Could not save.</p>`)
	}
	return render(c, http.StatusOK, templates.SaveOK())
}

func kindLabel(kind string) string {
	switch kind {
	case "fire":
		return "FIRE"
	case "sip":
		return "SIP"
	case "emi":
		return "EMI"
	case "emergency":
		return "Emergency"
	case "budget":
		return "Budget"
	default:
		return kind
	}
}
