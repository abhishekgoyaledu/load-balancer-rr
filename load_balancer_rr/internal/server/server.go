package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/load_balancer_rr/internal/config"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/backend"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/healthcheck"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/load_balancer"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/serverpool"
)

// Launch configuring the server and register all BE routes
func Launch() {
	// Create a root context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize a new server pool with lb algorithm
	serverPool, err := serverpool.NewServerPool()
	if err != nil {
		// Push alert here: Launch pool initialization failed
		config.Logger.Fatal(err.Error())
	}

	loadBalancer := load_balancer.NewLoadBalancer(serverPool)

	//executing for all services
	for _, routes := range config.Config.Backend.Routes {

		// Parse backend URLs and add them to the server pool
		parsedURL, parseErr := url.Parse(routes)
		if parseErr != nil {
			// Push alert here: URL parsing failed
			config.Logger.Fatal(parseErr.Error(), zap.String("URL", routes))
		}

		// Create a reverse proxy for the backend
		reverseProxy := httputil.NewSingleHostReverseProxy(parsedURL)

		// Create a new backend server and add it to the pool
		backendServer := backend.NewBackendServer(parsedURL, reverseProxy)

		serverPool.RegisterServiceBackend(backendServer)

		config.Logger.Info("added server", zap.String("host: ", backendServer.GetURL().Host))
	}
	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Server.Port),
		Handler:      http.HandlerFunc(loadBalancer.Serve),
		WriteTimeout: time.Duration(config.Config.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.Config.Server.ReadTimeout) * time.Second,
	}

	config.GracefulShutdownConfig(ctx, server)

	//running a go routing to perform healthcheck on the instances
	go healthcheck.PerformHealthCheck(ctx, serverPool)

	config.Logger.Info("Load Balancer is running successfully", zap.Int("port", config.Config.Server.Port))

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Push alert here: Launch encountered an unexpected error
		config.Logger.Fatal("Launch encountered an unexpected error", zap.Error(err))
	}
}
