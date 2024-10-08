package serverpool

import (
	"github.com/coda-payments/load_balancer_rr/internal/handlers/backend"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/serverpool/round_robin"
)

// ServerPool defines an interface for managing a pool of backend servers.
type ServerPool interface {
	// ListServiceBackends returns a slice of all backends in the pool.
	ListServiceBackends() []backend.Backend

	// NextAvailableBackend retrieves the next valid backend from the pool.
	NextAvailableBackend() backend.Backend

	// RegisterServiceBackend adds a new backend to the pool.
	RegisterServiceBackend(backend backend.Backend)

	// RemoveBackend the backend from the pool
	RemoveBackend(backend backend.Backend)

	// GetServerPoolSize returns the total number of backends in the pool.
	GetServerPoolSize() int32
}

// NewServerPool initializes and returns a new ServerPool instance.
// It creates a RoundRobin server pool with an empty list of backends.
func NewServerPool() (ServerPool, error) {
	return round_robin.Initialize(), nil
}
