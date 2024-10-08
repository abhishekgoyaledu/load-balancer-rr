package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/game_app/internal/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestLaunch tests the Launch function.
func TestLaunch(t *testing.T) {
	// Mock configurations
	config.Config.Server.WriteTimeout = 15
	config.Config.Server.ReadTimeout = 15

	// Set up a test server
	port := 9000
	serverAddr := fmt.Sprintf(":%d", port)
	v1Router := mux.NewRouter()
	healthcheckAPI(v1Router)
	routeAPIs(v1Router)
	server := &http.Server{
		Handler:      v1Router,
		Addr:         serverAddr,
		WriteTimeout: time.Duration(config.Config.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.Config.Server.ReadTimeout) * time.Second,
	}

	// Capture log output
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	config.Logger = logger

	// Run server in a goroutine to not block the test
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			config.Logger.Fatal("Launch encountered an unexpected error", zap.Error(err))
		}
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	// Check if the server is up by sending a request to the healthcheck endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthcheck", port))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Shutdown the server after test
	server.Close()
}

// TestHealthcheckAPI tests the healthcheckAPI function.
func TestHealthcheckAPI(t *testing.T) {
	router := mux.NewRouter()
	healthcheckAPI(router)

	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestRouteAPIs tests the routeAPIs function.
func TestRouteAPIs(t *testing.T) {
	router := mux.NewRouter()
	routeAPIs(router)

	// Prepare the game data to be sent in the request body
	gameData := map[string]interface{}{
		"game":    "Mobile Legends",
		"gamerID": "GYUTDTE",
		"points":  1000, // Assuming points is a variable that you define
	}

	// Convert the game data to JSON
	gameDataJSON, err := json.Marshal(gameData)
	assert.NoError(t, err, "Failed to marshal game data to JSON")

	// Create a new POST request with the game data JSON as the body
	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(gameDataJSON))
	assert.NoError(t, err, "Failed to create new request")

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a new ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Assuming CreatePlayer handler returns status 201 Created on success
	assert.Equal(t, http.StatusOK, rr.Code)
}
