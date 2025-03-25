package backend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestIsBackendAlive(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (*url.URL, func())
		expected bool
	}{
		{
			name: "backend is alive",
			setup: func() (*url.URL, func()) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				u, _ := url.Parse(server.URL)
				return u, server.Close
			},
			expected: true,
		},
		{
			name: "backend is not alive",
			setup: func() (*url.URL, func()) {
				// Use a non-existent port on localhost
				u, _ := url.Parse("http://localhost:12345")
				return u, func() {}
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, cleanup := tt.setup()
			defer cleanup()

			// Use a shorter timeout for the test
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			aliveChannel := make(chan bool, 1)
			go IsBackendAlive(ctx, aliveChannel, u)

			select {
			case result := <-aliveChannel:
				if result != tt.expected {
					t.Errorf("IsBackendAlive() = %v, want %v", result, tt.expected)
				}
			case <-ctx.Done():
				if tt.expected {
					t.Error("IsBackendAlive() timed out")
				}
			}
		})
	}
}
