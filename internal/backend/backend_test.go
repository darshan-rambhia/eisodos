package backend

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

func TestBackend_IsAlive(t *testing.T) {
	tests := []struct {
		name      string
		setAlive  bool
		wantAlive bool
	}{
		{
			name:      "backend is alive",
			setAlive:  true,
			wantAlive: true,
		},
		{
			name:      "backend is not alive",
			setAlive:  false,
			wantAlive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create backend
			serverURL, _ := url.Parse("http://test.com")
			b := NewBackend(serverURL, httputil.NewSingleHostReverseProxy(serverURL))
			backend := b.(*backend)

			// Initialize the channels
			backend.connections <- 0
			backend.SetAlive(tt.setAlive)

			// Test IsAlive
			got := backend.IsAlive()
			if got != tt.wantAlive {
				t.Errorf("Backend.IsAlive() = %v, want %v", got, tt.wantAlive)
			}
		})
	}
}

func TestBackend_GetURL(t *testing.T) {
	serverURL, _ := url.Parse("http://test.com")
	backend := NewBackend(serverURL, httputil.NewSingleHostReverseProxy(serverURL))
	if got := backend.GetURL(); got != serverURL {
		t.Errorf("Backend.GetURL() = %v, want %v", got, serverURL)
	}
}

func TestBackend_GetActiveConnections(t *testing.T) {
	serverURL, _ := url.Parse("http://test.com")
	b := NewBackend(serverURL, httputil.NewSingleHostReverseProxy(serverURL))
	backend := b.(*backend)

	// Initialize the connections channel
	backend.connections <- 0
	if got := backend.GetActiveConnections(); got != 0 {
		t.Errorf("Backend.GetActiveConnections() = %v, want 0", got)
	}
}

func TestBackend_ConnectionManagement(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	// Create backend
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %v", err)
	}
	b := NewBackend(serverURL, httputil.NewSingleHostReverseProxy(serverURL))
	backend := b.(*backend)

	// Initialize the channels
	backend.connections <- 0
	backend.SetAlive(true)

	// Test serving request
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	backend.Serve(recorder, req)

	if got := backend.GetActiveConnections(); got != 0 {
		t.Errorf("Backend.GetActiveConnections() = %v, want 0", got)
	}
}

func TestBackend_HealthCheck(t *testing.T) {
	// Create test server that responds after a delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create backend
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %v", err)
	}
	b := NewBackend(serverURL, httputil.NewSingleHostReverseProxy(serverURL))
	backend := b.(*backend)

	// Initialize the channels
	backend.connections <- 0
	backend.SetAlive(false)

	// Test health check
	if backend.IsAlive() {
		t.Error("Backend should be considered dead when marked as not alive")
	}

	backend.SetAlive(true)
	if !backend.IsAlive() {
		t.Error("Backend should be considered alive when marked as alive")
	}
}

func TestBackend_Serve(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	// Create backend
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %v", err)
	}
	b := NewBackend(serverURL, httputil.NewSingleHostReverseProxy(serverURL))
	backend := b.(*backend)

	// Initialize the channels
	backend.connections <- 0
	backend.SetAlive(true)

	// Test serving request
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	backend.Serve(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Backend.Serve() status code = %v, want %v", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "test response" {
		t.Errorf("Backend.Serve() body = %v, want %v", recorder.Body.String(), "test response")
	}
}
