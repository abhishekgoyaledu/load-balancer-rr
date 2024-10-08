package round_robin

import (
	"sync"
	"sync/atomic"

	"github.com/coda-payments/load_balancer_rr/internal/handlers/backend"
)

// RoundRobin represents a round-robin load balancer for backends.
type RoundRobin struct {
	Backends []backend.Backend
	Current  atomic.Int32
	mux      sync.RWMutex // Mutex for synchronizing access
}

// Initialize initializes and returns a new RoundRobin instance.
func Initialize() *RoundRobin {
	var rr RoundRobin
	rr.Backends = make([]backend.Backend, 0)
	rr.Current.Store(0)
	return &rr
}

// NextAvailableBackend returns the next alive backend in the pool.
func (roundRobin *RoundRobin) NextAvailableBackend() backend.Backend {
	var i int32
	for i = 0; i < roundRobin.GetServerPoolSize(); i++ {
		nextPeer := roundRobin.Rotate()
		alive := nextPeer.IsAlive()
		if alive.Load() {
			return nextPeer
		}
	}
	return nil
}

// Rotate moves to the next backend in the round-robin rotation.
func (roundRobin *RoundRobin) Rotate() backend.Backend {
	// Load the current index
	currentIndex := roundRobin.Current.Load()

	// Calculate the next index
	nextIndex := (currentIndex + 1) % roundRobin.GetServerPoolSize()

	// Store the updated index
	roundRobin.Current.Store(nextIndex)

	return roundRobin.Backends[nextIndex]
}

// ListServiceBackends returns all Backends in the pool.
func (roundRobin *RoundRobin) ListServiceBackends() []backend.Backend {
	roundRobin.mux.RLock()         // Acquire a read lock
	defer roundRobin.mux.RUnlock() // Ensure the read lock is released

	backendsCopy := make([]backend.Backend, len(roundRobin.Backends))
	copy(backendsCopy, roundRobin.Backends)
	return backendsCopy
}

// GetServerPoolSize returns the number of Backends in the pool.
func (roundRobin *RoundRobin) GetServerPoolSize() int32 {
	roundRobin.mux.RLock()
	defer roundRobin.mux.RUnlock()

	return int32(len(roundRobin.Backends))
}

// RegisterServiceBackend adds a new backend to the pool.
func (roundRobin *RoundRobin) RegisterServiceBackend(backend backend.Backend) {
	roundRobin.mux.Lock()         // Acquire a write lock
	defer roundRobin.mux.Unlock() //  Ensure the write lock is released

	roundRobin.Backends = append(roundRobin.Backends, backend)
}

// RemoveBackend removes a backend from the pool.
func (roundRobin *RoundRobin) RemoveBackend(backend backend.Backend) {
	roundRobin.mux.Lock()
	defer roundRobin.mux.Unlock()

	for i, b := range roundRobin.Backends {
		if b == backend {
			roundRobin.Backends = append(roundRobin.Backends[:i], roundRobin.Backends[i+1:]...)
			break
		}
	}
}
