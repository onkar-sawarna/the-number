package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNotFoundPage(t *testing.T) {
	e := echo.New()
	New().Register(e)
	req := httptest.NewRequest(http.MethodGet, "/no-such-page", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, `"message":"Not Found"`) {
		t.Fatal("JSON 404 leaked")
	}
	if !strings.Contains(body, "That page") || !strings.Contains(body, "href=\"/\"") {
		t.Fatalf("missing 404 copy: %s", body[:min(400, len(body))])
	}
}

func TestHomeStillOK(t *testing.T) {
	e := echo.New()
	New().Register(e)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "the number") {
		t.Fatal("home missing site name")
	}
}

func TestCrossingPage(t *testing.T) {
	e := echo.New()
	New().Register(e)
	req := httptest.NewRequest(http.MethodGet, "/crossing", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "The crossing") || !strings.Contains(body, "Jewellery") {
		t.Fatalf("crossing page missing copy: %s", body[:min(500, len(body))])
	}
}
