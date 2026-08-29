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
	if strings.Contains(rec.Body.String(), `href="/guidance"`) {
		t.Fatal("sleeves should not be a home chip")
	}
}

func TestFIREPageHasPlanLink(t *testing.T) {
	e := echo.New()
	New().Register(e)
	req := httptest.NewRequest(http.MethodGet, "/calculators/fire", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Copy this plan") {
		t.Fatal("FIRE page missing copy-this-plan")
	}
	if !strings.Contains(body, "If the year helped, UPI is welcome") {
		t.Fatal("FIRE page missing support line")
	}
	if !strings.Contains(body, "upi://pay?") {
		t.Fatal("FIRE page missing UPI pay link")
	}
	if !strings.Contains(body, "coastLine()") {
		t.Fatal("FIRE page missing coast line")
	}
	if strings.Contains(body, "Sleeves for your FIRE years") {
		t.Fatal("FIRE page should not push sleeves")
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
