package healthcheck

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/load_balancer_rr/internal/handlers/backend"
	"github.com/stretchr/testify/mock"

	"github.com/coda-payments/load_balancer_rr/internal/config"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/serverpool"
)

// MockServerPool is a mock implementation of the ServerPool interface.
type MockServerPool struct {
	mock.Mock
}

func (m *MockServerPool) ListServiceBackends() []backend.Backend {
	//TODO implement me
	panic("implement me")
}

func (m *MockServerPool) RemoveBackend(backend backend.Backend) {
	//TODO implement me
	panic("implement me")
}

func (m *MockServerPool) GetServerPoolSize() int32 {
	//TODO implement me
	panic("implement me")
}

func (m *MockServerPool) NextAvailableBackend() backend.Backend {
	args := m.Called()
	return args.Get(0).(backend.Backend)
}

func (m *MockServerPool) RegisterServiceBackend(backend backend.Backend) {
	m.Called(backend)
}

// MockHealthCheck is a mock function to simulate HealthCheck behavior.
func MockHealthCheck(ctx context.Context, sp serverpool.ServerPool) {
	// Simulate some health check behavior
}

// TestPerformHealthCheck tests the PerformHealthCheck function.
func TestPerformHealthCheck(t *testing.T) {
	// Set up mocks
	mockServerPool := new(MockServerPool)

	// Replace global Logger with a no-op logger for testing
	config.Logger, _ = zap.NewProduction()

	config.Config.HealthCheckTickerTimeInSeconds = int64(1)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Override the HealthCheck function with a mock
	originalHealthCheck := HealthCheck
	defer func() { HealthCheck = originalHealthCheck }()
	HealthCheck = MockHealthCheck

	// Start PerformHealthCheck in a separate goroutine
	go PerformHealthCheck(ctx, mockServerPool)

	// Allow some time for the ticker to trigger
	time.Sleep(3 * time.Second)

	// Cancel the context to stop health checks
	cancel()

	// Allow some time for context cancellation to be processed
	time.Sleep(1 * time.Second)
}
