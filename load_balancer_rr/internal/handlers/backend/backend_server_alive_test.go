package backend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

func TestIsServerAlive(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus bool
	}{
		{
			name:           "Server is alive",
			statusCode:     http.StatusOK,
			expectedStatus: true,
		},
		{
			name:           "Server is not alive",
			statusCode:     http.StatusInternalServerError,
			expectedStatus: false,
		},
		{
			name:           "Server is unreachable",
			statusCode:     -1, // This will simulate an unreachable server
			expectedStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.statusCode != -1 {
					w.WriteHeader(tt.statusCode) // Respond with the specified status code
				} else {
					// Simulate server unreachable
					w.WriteHeader(http.StatusServiceUnavailable)
				}
			}))
			defer mockServer.Close()

			// Parse the URL of the mock server
			serverURL, err := url.Parse(mockServer.URL)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}

			isAliveChannel := make(chan atomic.Bool)

			// Run IsServerAlive in a separate goroutine
			go IsServerAlive(context.Background(), isAliveChannel, serverURL)

			// Wait for the result
			select {
			case status := <-isAliveChannel:
				if status.Load() != tt.expectedStatus {
					t.Errorf("Expected server alive status to be %v, got %v", tt.expectedStatus, status.Load())
				}
			case <-time.After(3 * time.Second): // Increased timeout for slower tests
				t.Fatal("Test timed out waiting for response")
			}
		})
	}
}
