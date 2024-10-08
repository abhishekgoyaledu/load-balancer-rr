package load_balancer

import (
	"net/http"

	"github.com/coda-payments/load_balancer_rr/internal/handlers/serverpool"
)

// LoadBalancer defines the interface for a load balancer that can serve HTTP requests.
type LoadBalancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

// loadBalancer is a concrete implementation of the LoadBalancer interface.
type loadBalancer struct {
	serverPool serverpool.ServerPool
}

// Serve handles incoming HTTP requests by forwarding them to the next available backend server.
func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	// Get the next available backend server from the server pool.
	backend := lb.serverPool.NextAvailableBackend()
	if backend != nil {
		// If a backend server is available, forward the request to it.
		backend.Serve(w, r)
		return
	}
	// If no backend server is available, respond with a 503 Service Unavailable error.
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// NewLoadBalancer creates a new instance of a load balancer with the specified server pool.
func NewLoadBalancer(serverPool serverpool.ServerPool) LoadBalancer {
	return &loadBalancer{
		serverPool: serverPool,
	}
}
