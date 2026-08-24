package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thenumber/app/internal/models"
	"github.com/thenumber/app/internal/session"
	"github.com/thenumber/app/web/templates"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) loginGet(c echo.Context) error {
	next := session.SafeNext(c.QueryParam("next"))
	return render(c, http.StatusOK, templates.Login(templates.AuthView{
		Page: s.page(c, "Log in"),
		Next: next,
	}))
}

func (s *Server) loginPost(c echo.Context) error {
	email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	password := c.FormValue("password")
	next := session.SafeNext(c.FormValue("next"))
	if next == "" {
		next = session.SafeNext(c.QueryParam("next"))
	}
	view := templates.AuthView{Page: s.page(c, "Log in"), Value: email, Next: next}

	var u models.User
	err := s.DB.Where("email = ?", email).First(&u).Error
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		view.Error = "Email or password is incorrect."
		return render(c, http.StatusUnauthorized, templates.Login(view))
	}
	if err := session.SaveUserID(c, u.ID); err != nil {
		view.Error = "Could not start a session."
		return render(c, http.StatusInternalServerError, templates.Login(view))
	}
	if next == "" {
		next = "/dashboard"
	}
	return c.Redirect(http.StatusSeeOther, next)
}

func (s *Server) signupGet(c echo.Context) error {
	return render(c, http.StatusOK, templates.Signup(templates.AuthView{Page: s.page(c, "Sign up")}))
}

func (s *Server) signupPost(c echo.Context) error {
	email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	password := c.FormValue("password")
	confirm := c.FormValue("confirm_password")
	view := templates.AuthView{Page: s.page(c, "Sign up"), Value: email}

	if !strings.Contains(email, "@") || len(email) < 5 {
		view.Error = "Enter a valid email."
		return render(c, http.StatusUnprocessableEntity, templates.Signup(view))
	}
	if len(password) < 8 {
		view.Error = "Password must be at least 8 characters."
		return render(c, http.StatusUnprocessableEntity, templates.Signup(view))
	}
	if password != confirm {
		view.Error = "Passwords do not match."
		return render(c, http.StatusUnprocessableEntity, templates.Signup(view))
	}

	var n int64
	if err := s.DB.Model(&models.User{}).Where("email = ?", email).Count(&n).Error; err != nil {
		view.Error = "Could not create account."
		return render(c, http.StatusInternalServerError, templates.Signup(view))
	}
	if n > 0 {
		view.Error = "An account with this email already exists."
		return render(c, http.StatusConflict, templates.Signup(view))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		view.Error = "Could not create account."
		return render(c, http.StatusInternalServerError, templates.Signup(view))
	}
	u := models.User{Email: email, PasswordHash: string(hash)}
	if err := s.DB.Create(&u).Error; err != nil {
		view.Error = "An account with this email already exists."
		return render(c, http.StatusConflict, templates.Signup(view))
	}
	if err := session.SaveUserID(c, u.ID); err != nil {
		view.Error = "Account created; please log in."
		return render(c, http.StatusInternalServerError, templates.Signup(view))
	}
	return c.Redirect(http.StatusSeeOther, "/calculators/fire")
}

func (s *Server) logout(c echo.Context) error {
	_ = session.Clear(c)
	return c.Redirect(http.StatusSeeOther, "/")
}
