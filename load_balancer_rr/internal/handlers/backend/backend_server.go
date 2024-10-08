package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// backendServer represents a single backendServer server with its URL and state.
type backendServer struct {
	url          *url.URL
	alive        atomic.Bool
	reverseProxy *httputil.ReverseProxy
}

// Backend interface defines methods for interacting with a backendServer server.
type Backend interface {
	Serve(http.ResponseWriter, *http.Request)
	GetURL() *url.URL
	SetAlive(atomic.Bool)
	IsAlive() atomic.Bool
}

// NewBackendServer initializes and returns a new backendServer instance.
func NewBackendServer(u *url.URL, rp *httputil.ReverseProxy) Backend {
	server := &backendServer{
		url:          u,
		reverseProxy: rp,
	}
	server.alive.Store(true)
	return server
}

// SetAlive updates the alive state of the backendServer server.
func (b *backendServer) SetAlive(alive atomic.Bool) {
	b.alive = alive
}

// IsAlive checks if the backendServer server is alive.
func (b *backendServer) IsAlive() atomic.Bool {
	return b.alive
}

// GetURL retrieves the URL of the backendServer server.
func (b *backendServer) GetURL() *url.URL {
	return b.url
}

// Serve handles incoming HTTP requests and forwards them to the backendServer server.
func (b *backendServer) Serve(rw http.ResponseWriter, req *http.Request) {
	//push an alert here to check how many request we are triggering to each instance
	// Proxy the request to the backendServer server
	b.reverseProxy.ServeHTTP(rw, req)
}
