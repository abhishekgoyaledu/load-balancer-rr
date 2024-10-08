package player

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Successful Request",
			payload:        map[string]interface{}{"key1": "value1", "key2": 123.0},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"key1": "value1", "key2": 123.0},
		},
		{
			name:           "Invalid JSON",
			payload:        nil, // Will not be used in this case
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if tt.name == "Successful Request" {
				reqBody, _ = json.Marshal(tt.payload)
			} else {
				reqBody = []byte(`invalid json`) // Invalid JSON for testing
			}

			req := httptest.NewRequest(http.MethodPost, "/player", bytes.NewBuffer(reqBody))
			rec := httptest.NewRecorder()

			CreatePlayer(rec, req)

			res := rec.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var responseBody map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&responseBody)
				if err != nil {
					t.Errorf("Failed to decode response body: %v", err)
				}
				if !equalMaps(responseBody, tt.expectedBody) {
					t.Errorf("expected body %+v, got %+v", tt.expectedBody, responseBody)
				}
			}
		})
	}
}

// Helper function to compare two maps
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}
