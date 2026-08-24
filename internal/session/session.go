package session

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const CookieName = "the-number-session"

// DefaultSecret is a development-only signing key. Change SESSION_SECRET in any shared environment.
const DefaultSecret = "dev-only-the-number-session-secret-change-me"

func Secret() string {
	if s := os.Getenv("SESSION_SECRET"); s != "" {
		return s
	}
	return DefaultSecret
}

func NewStore() *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(Secret()))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   14 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return store
}

func Get(c echo.Context) (*sessions.Session, error) {
	return session.Get(CookieName, c)
}

func SaveUserID(c echo.Context, userID uint) error {
	sess, err := Get(c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   14 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["user_id"] = userID
	return sess.Save(c.Request(), c.Response())
}

func Clear(c echo.Context) error {
	sess, err := Get(c)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	delete(sess.Values, "user_id")
	return sess.Save(c.Request(), c.Response())
}

func UserID(c echo.Context) (uint, bool) {
	sess, err := Get(c)
	if err != nil {
		return 0, false
	}
	switch v := sess.Values["user_id"].(type) {
	case uint:
		return v, v != 0
	case int:
		return uint(v), v != 0
	case int64:
		return uint(v), v != 0
	default:
		return 0, false
	}
}

// SafeNext allows only relative paths (no protocol, no scheme-relative URLs).
func SafeNext(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || !strings.HasPrefix(s, "/") || strings.HasPrefix(s, "//") || strings.Contains(s, "://") {
		return ""
	}
	if strings.ContainsAny(s, "\r\n") {
		return ""
	}
	return s
}
