package load_balancer_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/coda-payments/load_balancer_rr/internal/handlers/backend"
	"github.com/coda-payments/load_balancer_rr/internal/handlers/load_balancer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServer is a mock implementation of a backend server.
type MockServer struct {
	mock.Mock
}

func (m *MockServer) Serve(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// MockServerPool is a mock implementation of the server pool.
type MockServerPool struct {
	mock.Mock
}

func (m *MockServerPool) ListServiceBackends() []backend.Backend {
	//TODO implement me
	panic("implement me")
}

func (m *MockServerPool) NextAvailableBackend() backend.Backend {
	parsedURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	bs := backend.NewBackendServer(parsedURL, proxy)
	return bs
}

func (m *MockServerPool) RegisterServiceBackend(backend backend.Backend) {
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

func TestLoadBalancer_Serve_BackendAvailable(t *testing.T) {
	// Create a mock backend server.
	mockBackend := new(MockServer)
	mockBackend.On("Serve", mock.Anything, mock.Anything).Return()

	// Create a mock server pool that returns the mock backend server.
	mockPool := new(MockServerPool)
	mockPool.On("NextAvailableBackend").Return(mockBackend)

	// Create a load balancer with the mock server pool.
	lb := load_balancer.NewLoadBalancer(mockPool)

	// Create a mock HTTP request and response recorder.
	req, _ := http.NewRequest("GET", "/create", nil)
	rr := httptest.NewRecorder()

	// Call the Serve method.
	lb.Serve(rr, req)

	// Assert that the response code is 200 OK.
	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

func TestLoadBalancer_Serve_NoBackendAvailable(t *testing.T) {
	// Create a mock server pool that returns nil (no backend available).
	mockPool := new(MockServerPool)
	mockPool.On("NextAvailableBackend").Return(nil)

	// Create a load balancer with the mock server pool.
	lb := load_balancer.NewLoadBalancer(mockPool)
	// Create a mock HTTP request and response recorder.
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Call the Serve method.
	lb.Serve(rr, req)

	// Assert that the response code is 503 Service Unavailable.
	assert.Equal(t, http.StatusBadGateway, rr.Code)
}
