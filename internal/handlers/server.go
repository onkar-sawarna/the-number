package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/thenumber/app/internal/calc"
	"github.com/thenumber/app/web/templates"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Register(e *echo.Echo) {
	e.GET("/", s.landing)
	e.GET("/about", s.about)
	e.GET("/crossing", s.crossing)
	e.GET("/disclaimer", s.disclaimer)
	e.POST("/theme", s.theme)

	e.GET("/calculators/fire", s.firePage)
	e.GET("/calculators/sip", s.sipPage)
	e.GET("/calculators/emi", s.emiPage)
	e.GET("/calculators/emergency", s.emergencyPage)
	e.GET("/calculators/budget", s.budgetPage)

	e.GET("/guidance", s.guidanceGet)
	e.POST("/guidance", s.guidancePost)

	fallback := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		s.httpError(err, c, fallback)
	}
}

func (s *Server) httpError(err error, c echo.Context, fallback echo.HTTPErrorHandler) {
	code := http.StatusInternalServerError
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
	}
	if c.Response().Committed {
		return
	}
	if code == http.StatusNotFound {
		if renderErr := render(c, http.StatusNotFound, templates.NotFound(s.page(c, "Not found", "title_404"))); renderErr != nil {
			c.Logger().Error(renderErr)
		}
		return
	}
	fallback(err, c)
}

func render(c echo.Context, status int, t templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(status)
	return t.Render(c.Request().Context(), c.Response())
}

func (s *Server) page(c echo.Context, title, titleKey string) templates.Page {
	return templates.Page{
		Title:    title,
		TitleKey: titleKey,
		Dark:     themeDark(c),
		Path:     c.Request().URL.Path,
	}
}

func themeDark(c echo.Context) bool {
	ck, err := c.Cookie("theme")
	return err == nil && ck.Value == "dark"
}

func (s *Server) landing(c echo.Context) error {
	return render(c, http.StatusOK, templates.Landing(s.page(c, "Home", "title_home")))
}

func (s *Server) about(c echo.Context) error {
	return render(c, http.StatusOK, templates.About(s.page(c, "About", "title_about")))
}

func (s *Server) crossing(c echo.Context) error {
	p := s.page(c, "The crossing", "title_crossing")
	p.Description = templates.Copy("crossing_meta")
	return render(c, http.StatusOK, templates.Crossing(p))
}

func (s *Server) disclaimer(c echo.Context) error {
	return render(c, http.StatusOK, templates.Disclaimer(s.page(c, "Disclaimer", "title_disclaimer")))
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
	s := strings.TrimSpace(c.FormValue(name))
	if s == "" {
		return def
	}
	v, ok := calc.ParseCompactMoney(s)
	if !ok {
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
