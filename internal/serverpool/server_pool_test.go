package serverpool

import (
	"context"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

func TestNewServerPool(t *testing.T) {
	tests := []struct {
		name     string
		strategy LBStrategy
		wantErr  bool
	}{
		{
			name:     "round-robin strategy",
			strategy: RoundRobin,
			wantErr:  false,
		},
		{
			name:     "least-connected strategy",
			strategy: LeastConnected,
			wantErr:  false,
		},
		{
			name:     "invalid strategy",
			strategy: LBStrategy(999), // Invalid strategy
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, err := NewServerPool(tt.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && sp == nil {
				t.Error("NewServerPool() returned nil server pool")
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	sp, err := NewServerPool(RoundRobin)
	if err != nil {
		t.Fatalf("Failed to create server pool: %v", err)
	}

	// Add a mock backend
	urlStr := "http://localhost:8080"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	mb := newMockBackend(u, proxy)
	sp.AddBackend(mb)

	// Test health check
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		HealthCheck(ctx, sp)
		close(done)
	}()

	// Wait for at least one health check
	time.Sleep(150 * time.Millisecond)

	// Cancel context and wait for health check to stop
	cancel()
	select {
	case <-done:
		// Health check stopped as expected
	case <-time.After(1 * time.Second):
		t.Error("HealthCheck did not stop after context cancellation")
	}
}
