package serverpool

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/darshan-rambhia/eisodos/internal/backend"
)

func TestLeastConnectedServerPool(t *testing.T) {
	pool := &lcServerPool{
		backends: make([]backend.Backend, 0),
	}

	// Test empty pool
	if got := pool.GetNextValidPeer(); got != nil {
		t.Errorf("GetNextValidPeer() = %v, want nil", got)
	}

	// Create test backends
	urls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, u := range urls {
		url, _ := url.Parse(u)
		proxy := httputil.NewSingleHostReverseProxy(url)
		backend := newMockBackend(url, proxy)
		pool.AddBackend(backend)
	}

	// Test pool size
	if got := pool.GetServerPoolSize(); got != len(urls) {
		t.Errorf("GetServerPoolSize() = %v, want %v", got, len(urls))
	}

	// Test least-connected behavior
	// All backends should have 0 connections initially
	backend1 := pool.GetNextValidPeer()
	if backend1 == nil {
		t.Fatal("GetNextValidPeer() = nil, want non-nil")
	}

	// Simulate some connections
	for i := 0; i < 3; i++ {
		backend1.Serve(nil, nil) // This will increment the connection count
	}

	// Get next backend - should be a different one with fewer connections
	backend2 := pool.GetNextValidPeer()
	if backend2 == nil {
		t.Fatal("GetNextValidPeer() = nil, want non-nil")
	}

	if backend2 == backend1 {
		t.Error("GetNextValidPeer() returned same backend with more connections")
	}

	// Test GetBackends
	backends := pool.GetBackends()
	if len(backends) != len(urls) {
		t.Errorf("GetBackends() returned %v backends, want %v", len(backends), len(urls))
	}
}
