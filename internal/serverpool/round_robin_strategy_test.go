package serverpool

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/darshan-rambhia/eisodos/internal/backend"
)

func TestRoundRobinServerPool(t *testing.T) {
	pool := &roundRobinServerPool{
		backends: make([]backend.Backend, 0),
		current:  0,
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

	// Test round-robin behavior
	var lastBackend backend.Backend
	for i := 0; i < len(urls)*2; i++ {
		backend := pool.GetNextValidPeer()
		if backend == nil {
			t.Errorf("GetNextValidPeer() = nil, want non-nil")
			continue
		}

		if i >= len(urls) && backend == lastBackend {
			t.Errorf("GetNextValidPeer() returned same backend twice in a row")
		}
		lastBackend = backend
	}

	// Test GetBackends
	backends := pool.GetBackends()
	if len(backends) != len(urls) {
		t.Errorf("GetBackends() returned %v backends, want %v", len(backends), len(urls))
	}
}
