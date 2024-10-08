package backend

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"testing"
)

// TestNewBackendServer tests the initialization of a backendServer instance.
func TestNewBackendServer(t *testing.T) {
	parsedURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	bs := NewBackendServer(parsedURL, proxy)

	if bs.GetURL().String() != "http://localhost:8080" {
		t.Errorf("Expected URL to be 'http://localhost:8080', got '%s'", bs.GetURL())
	}

	alive := bs.IsAlive()
	if !alive.Load() {
		t.Error("Expected backendServer to be alive upon initialization")
	}
}

// TestSetAlive tests the SetAlive method.
func TestSetAlive(t *testing.T) {
	aliveStatusFalse := atomic.Bool{}
	aliveStatusTrue := atomic.Bool{}

	aliveStatusFalse.Store(false)
	aliveStatusTrue.Store(true)
	parsedURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	bs := NewBackendServer(parsedURL, proxy)

	bs.SetAlive(aliveStatusFalse)
	alive := bs.IsAlive()
	if alive.Load() {
		t.Error("Expected backendServer to be not alive after setting to false")
	}

	bs.SetAlive(aliveStatusTrue)
	isAlive := bs.IsAlive()
	if !isAlive.Load() {
		t.Error("Expected backendServer to be alive after setting to true")
	}
}

// TestGetURL tests the GetURL method.
func TestGetURL(t *testing.T) {
	parsedURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	bs := NewBackendServer(parsedURL, proxy)

	if bs.GetURL() != parsedURL {
		t.Error("GetURL did not return the expected URL")
	}
}
