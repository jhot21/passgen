package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	web "codeberg.org/jhot/passgen/web"
)

func newServer() http.Handler { return web.NewServer() }

// parseBool tests

func TestParseBool(t *testing.T) {
	truthy := []string{"true", "True", "TRUE", "yes", "Yes", "YES", "1", "on", "On", "ON"}
	falsy := []string{"false", "False", "no", "0", "off", "", "maybe", "nope"}

	for _, s := range truthy {
		if !web.ParseBool(s) {
			t.Errorf("ParseBool(%q) = false, want true", s)
		}
	}
	for _, s := range falsy {
		if web.ParseBool(s) {
			t.Errorf("ParseBool(%q) = true, want false", s)
		}
	}
}

// GET tests

func TestGET_PlainText_Default(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/generate?passphrase=testphrase", nil)
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("want text/plain Content-Type, got %q", ct)
	}
	body := strings.TrimSpace(w.Body.String())
	if len(body) == 0 {
		t.Fatal("expected non-empty password")
	}
	if strings.Contains(body, "{") {
		t.Fatalf("expected plain text, got JSON: %s", body)
	}
}

func TestGET_JSON_AcceptHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/generate?passphrase=testphrase", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if resp["password"] == "" {
		t.Fatal("expected non-empty password in JSON response")
	}
}

func TestGET_MissingPassphrase_Returns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/generate", nil)
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestGET_Deterministic(t *testing.T) {
	run := func() string {
		req := httptest.NewRequest(http.MethodGet, "/api/generate?passphrase=hello&salt=world&length=20", nil)
		w := httptest.NewRecorder()
		newServer().ServeHTTP(w, req)
		return strings.TrimSpace(w.Body.String())
	}
	a, b := run(), run()
	if a != b {
		t.Errorf("expected deterministic output: %q != %q", a, b)
	}
}

func TestGET_BoolParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/generate?passphrase=test&no_specials=1&no_uppers=yes&length=40", nil)
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	body := strings.TrimSpace(w.Body.String())
	for _, c := range body {
		if strings.ContainsRune("!@#$%^&*", c) {
			t.Errorf("password contains special char %q but no_specials=1 was set", c)
		}
		if c >= 'A' && c <= 'Z' {
			t.Errorf("password contains uppercase %q but no_uppers=yes was set", c)
		}
	}
}

// POST tests

func TestPOST_JSON(t *testing.T) {
	body := `{"passphrase":"testphrase","salt":"testsalt","length":24}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(resp["password"]) != 24 {
		t.Errorf("want password of length 24, got %d: %q", len(resp["password"]), resp["password"])
	}
}

func TestPOST_InvalidJSON_Returns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestPOST_PlainTextResponse(t *testing.T) {
	body := `{"passphrase":"testphrase"}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No Accept header → plain text
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("want text/plain, got %q", ct)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/generate", nil)
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", w.Code)
	}
}

func TestGET_CustomSpecials(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/generate?passphrase=test&custom_specials=!&length=40", nil)
	w := httptest.NewRecorder()
	newServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	body := strings.TrimSpace(w.Body.String())
	// With custom_specials=! only '!' is allowed as a special char
	// Verify none of the other default specials appear
	for _, c := range body {
		if strings.ContainsRune("@#$%^&*", c) {
			t.Errorf("password contains %q but custom_specials=! was set", c)
		}
	}
}
