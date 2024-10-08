package round_robin

import (
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockBackend is a mock implementation of the backend.Backend interface
type MockBackend struct {
	alive   atomic.Bool
	address string
	url     *url.URL
}

// IsAlive returns the alive status of the backend.
func (m *MockBackend) IsAlive() atomic.Bool {
	return m.alive
}

// SetAlive sets the alive status of the backend.
func (m *MockBackend) SetAlive(alive atomic.Bool) {
	m.alive.Store(alive.Load())
}

// GetAddress returns the address of the backend.
func (m *MockBackend) GetAddress() string {
	return m.address
}

// SetAddress sets the address of the backend.
func (m *MockBackend) SetAddress(address string) {
	m.address = address
}

// Serve mocks the serving of HTTP requests.
func (m *MockBackend) Serve(w http.ResponseWriter, r *http.Request) {
	if !m.alive.Load() {
		http.Error(w, "Backend is down", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Response from backend at " + m.address))
}

// GetURL returns the URL of the backend.
func (m *MockBackend) GetURL() *url.URL {
	if m.url == nil {
		parsedURL, _ := url.Parse(m.address)
		m.url = parsedURL
	}
	return m.url
}

func TestNew(t *testing.T) {
	rr := Initialize()
	assert.NotNil(t, rr)
	assert.Equal(t, int32(0), rr.Current.Load())
	assert.Empty(t, rr.Backends)
}

func TestRegisterServiceBackend(t *testing.T) {
	rr := Initialize()
	mockBackend := &MockBackend{}
	rr.RegisterServiceBackend(mockBackend)

	assert.Len(t, rr.Backends, 1)
	assert.Equal(t, mockBackend, rr.Backends[0])
}

func TestRemoveBackend(t *testing.T) {
	rr := Initialize()
	mockBackend1 := &MockBackend{}
	mockBackend2 := &MockBackend{}
	rr.RegisterServiceBackend(mockBackend1)
	rr.RegisterServiceBackend(mockBackend2)

	rr.RemoveBackend(mockBackend1)
	assert.Len(t, rr.Backends, 1)
	assert.Equal(t, mockBackend2, rr.Backends[0])

	rr.RemoveBackend(mockBackend2)
	assert.Empty(t, rr.Backends)
}

func TestRotate(t *testing.T) {
	rr := Initialize()
	mockBackend1 := &MockBackend{}
	mockBackend2 := &MockBackend{}
	rr.RegisterServiceBackend(mockBackend1)
	rr.RegisterServiceBackend(mockBackend2)

	rr.Current.Store(0)
	assert.Equal(t, mockBackend2, rr.Rotate())
	assert.Equal(t, int32(1), rr.Current.Load())

	assert.Equal(t, mockBackend1, rr.Rotate())
	assert.Equal(t, int32(0), rr.Current.Load())
}

func TestNextAvailableBackend(t *testing.T) {
	rr := Initialize()
	aliveStatusFalse := atomic.Bool{}
	aliveStatusTrue := atomic.Bool{}

	aliveStatusFalse.Store(false)
	aliveStatusTrue.Store(true)

	mockBackend1 := &MockBackend{}
	mockBackend2 := &MockBackend{}
	mockBackend3 := &MockBackend{}
	mockBackend1.SetAlive(aliveStatusFalse)

	mockBackend2.SetAlive(aliveStatusTrue)
	mockBackend3.SetAlive(aliveStatusTrue)
	rr.RegisterServiceBackend(mockBackend1)
	rr.RegisterServiceBackend(mockBackend2)
	rr.RegisterServiceBackend(mockBackend3)

	backend := rr.NextAvailableBackend()
	assert.NotNil(t, backend)
	assert.Equal(t, mockBackend2, backend)

	mockBackend2.SetAlive(aliveStatusFalse)
	backend = rr.NextAvailableBackend()
	assert.Equal(t, mockBackend3, backend)

	mockBackend3.SetAlive(aliveStatusFalse)
	backend = rr.NextAvailableBackend()
	assert.Nil(t, backend)
}

func TestListServiceBackends(t *testing.T) {
	rr := Initialize()
	mockBackend1 := &MockBackend{}
	mockBackend2 := &MockBackend{}
	rr.RegisterServiceBackend(mockBackend1)
	rr.RegisterServiceBackend(mockBackend2)

	backends := rr.ListServiceBackends()
	assert.Len(t, backends, 2)
	assert.Equal(t, mockBackend1, backends[0])
	assert.Equal(t, mockBackend2, backends[1])
}

func TestGetServerPoolSize(t *testing.T) {
	rr := Initialize()
	assert.Equal(t, int32(0), rr.GetServerPoolSize())

	mockBackend1 := &MockBackend{}
	rr.RegisterServiceBackend(mockBackend1)
	assert.Equal(t, int32(1), rr.GetServerPoolSize())

	mockBackend2 := &MockBackend{}
	rr.RegisterServiceBackend(mockBackend2)
	assert.Equal(t, int32(2), rr.GetServerPoolSize())
}
