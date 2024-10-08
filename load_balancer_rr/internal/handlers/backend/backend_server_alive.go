package backend

import (
	"context"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/coda-payments/load_balancer_rr/internal/config"
	"github.com/coda-payments/load_balancer_rr/internal/constant"
)

// IsServerAlive checks if the server is alive by sending an HTTP GET request
// to the server health check endpoint.
var IsServerAlive = func(ctx context.Context, isAliveChannel chan atomic.Bool, url *url.URL) {
	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: time.Duration(config.Config.Backend.Endpoint[constant.Healthcheck].Timeout) * time.Second, // Set a timeout for the request
	}

	urlString := url.String() + config.Config.Backend.Endpoint[constant.Healthcheck].URL

	// Create the HTTP request to the health check endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		//push a metrics
		isAliveChannel <- getAliveStatus(false)
		return
	}

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		// If there's an error, the server is not alive
		//push a metrics
		isAliveChannel <- getAliveStatus(false)
		return
	}
	defer resp.Body.Close()

	// Check the status code received from healthcheck API to determine if the server is alive
	var aliveStatus bool
	if resp.StatusCode == http.StatusOK {
		aliveStatus = true
	} else {
		aliveStatus = false
	}
	isAliveChannel <- getAliveStatus(aliveStatus)
}

// getAliveStatus returns an atomic.Bool representing the alive status of the server
func getAliveStatus(status bool) atomic.Bool {
	aliveStatus := atomic.Bool{}
	aliveStatus.Store(status)
	return aliveStatus
}
