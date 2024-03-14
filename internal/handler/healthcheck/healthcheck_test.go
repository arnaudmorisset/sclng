package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPongHandler(t *testing.T) {
	handler := NewPongHandler()

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler(rr, req, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	expectedBody := `{"status":"pong"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rr.Body.String())
	}
}
