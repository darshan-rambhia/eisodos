package main

import (
	"testing"
	"time"

	"github.com/darshan-rambhia/eisodos/internal/serverpool"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromCLI(t *testing.T) {
	tests := []struct {
		name                string
		port                int
		healthCheckInterval time.Duration
		strategy            string
		backendURLs         []string
		wantErr             bool
		errContains         string
	}{
		{
			name:                "valid configuration",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "round-robin",
			backendURLs:         []string{"http://localhost:8081", "http://localhost:8082"},
			wantErr:             false,
		},
		{
			name:                "no backends",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "round-robin",
			backendURLs:         []string{},
			wantErr:             true,
			errContains:         "please provide at least one backend URL",
		},
		{
			name:                "invalid backend URL",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "round-robin",
			backendURLs:         []string{"invalid-url"},
			wantErr:             true,
			errContains:         "failed to parse backend URL",
		},
		{
			name:                "invalid strategy",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "invalid-strategy",
			backendURLs:         []string{"http://localhost:8081"},
			wantErr:             false, // Should default to round-robin
		},
		{
			name:                "backend URL with missing scheme",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "round-robin",
			backendURLs:         []string{"localhost:8081"},
			wantErr:             true,
			errContains:         "missing scheme or host",
		},
		{
			name:                "backend URL with missing host",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "round-robin",
			backendURLs:         []string{"http://"},
			wantErr:             true,
			errContains:         "missing scheme or host",
		},
		{
			name:                "multiple backends with one invalid",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "round-robin",
			backendURLs:         []string{"http://localhost:8081", "invalid-url", "http://localhost:8083"},
			wantErr:             true,
			errContains:         "failed to parse backend URL",
		},
		{
			name:                "least-connected strategy",
			port:                8080,
			healthCheckInterval: 10 * time.Second,
			strategy:            "least-connected",
			backendURLs:         []string{"http://localhost:8081", "http://localhost:8082"},
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb, err := LoadFromCLI(tt.port, tt.healthCheckInterval, tt.strategy, tt.backendURLs)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, lb)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, lb)

			// Verify the load balancer configuration
			assert.Equal(t, tt.port, lb.GetPort())
			assert.Equal(t, len(tt.backendURLs), lb.GetServerPoolSize())

			// Verify the backends
			backends := lb.GetBackends()
			assert.Len(t, backends, len(tt.backendURLs))
			for i, b := range backends {
				assert.Equal(t, tt.backendURLs[i], b.GetURL().String())
			}
		})
	}
}

func TestParseStrategy(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
		want     serverpool.LBStrategy
	}{
		{
			name:     "round-robin strategy",
			strategy: "round-robin",
			want:     serverpool.RoundRobin,
		},
		{
			name:     "least-connected strategy",
			strategy: "least-connected",
			want:     serverpool.LeastConnected,
		},
		{
			name:     "unknown strategy defaults to round-robin",
			strategy: "unknown",
			want:     serverpool.RoundRobin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serverpool.ParseStrategy(tt.strategy)
			assert.Equal(t, tt.want, got)
		})
	}
}
