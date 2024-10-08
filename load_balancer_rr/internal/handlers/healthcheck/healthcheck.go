package healthcheck

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/load_balancer_rr/internal/config"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/backend"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/serverpool"
)

const (
	HealthyStatus   = "Healthy"
	UnhealthyStatus = "Unhealthy"
)

// PerformHealthCheck initiates a periodic health check for backend hosts in the server pool.
func PerformHealthCheck(ctx context.Context, sp serverpool.ServerPool) {
	config.Logger.Info("Starting health check for backend hosts")
	// Create a ticker to perform health checks at specified intervals.
	ticker := time.NewTicker(time.Duration(config.Config.HealthCheckTickerTimeInSeconds) * time.Second)

	for {
		select {
		// Trigger health check on each tick.
		case <-ticker.C:
			go HealthCheck(ctx, sp)
		// Handle context cancellation to gracefully stop health check execution.
		case <-ctx.Done():
			config.Logger.Info("Closing health check execution..")
			return
		}
	}
}

// HealthCheck verifies the status of each service backend in the server pool.
var HealthCheck = func(ctx context.Context, sp serverpool.ServerPool) {
	aliveChannel := make(chan atomic.Bool, 1) // Channel to receive the alive status of each service.

	for _, service := range sp.ListServiceBackends() {
		// Create a new context with a timeout for the health check request.
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		healthStatus := HealthyStatus

		// Asynchronously check if the backend service is alive.
		go backend.IsServerAlive(requestCtx, aliveChannel, service.GetURL())

		select {
		// Handle context cancellation, logging a shutdown message.
		case <-ctx.Done():
			config.Logger.Info("Gracefully shutting down health check")
			return
		// Wait for the alive status from the channel.
		case alive := <-aliveChannel:
			service.SetAlive(alive)
			if !alive.Load() {
				// Push an alert here for a health check failure or configure the number of hosts.
				healthStatus = UnhealthyStatus
			}
		}
		// Call stop to cancel the context for the current service health check.
		stop()

		config.Logger.Info("host status: ", zap.String("status", healthStatus), zap.String("host", service.GetURL().String()))
	}
}
