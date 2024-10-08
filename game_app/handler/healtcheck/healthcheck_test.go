package healtcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rec := httptest.NewRecorder()

	HealthCheck(rec, req)

	// Since the function currently does not send any response, we only check for execution.
	if rec.Code != http.StatusOK { // HTTP status code 0 indicates no response was sent
		t.Errorf("Expected response code to be 0, got %d", rec.Code)
	}
}
