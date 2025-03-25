package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/darshan-rambhia/eisodos/internal/serverpool"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 10*time.Second, cfg.HealthCheckInterval)
	assert.Equal(t, serverpool.RoundRobin, cfg.Strategy)
	assert.Empty(t, cfg.Backends)
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid configuration",
			config: &Config{
				Port:                8080,
				HealthCheckInterval: 10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends: []BackendConfig{
					{
						URL:      "http://localhost:8081",
						Weight:   1,
						MaxConns: 100,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: &Config{
				Port:                -1,
				HealthCheckInterval: 10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends: []BackendConfig{
					{URL: "http://localhost:8081"},
				},
			},
			wantErr:     true,
			errContains: "invalid port number",
		},
		{
			name: "invalid health check interval",
			config: &Config{
				Port:                8080,
				HealthCheckInterval: -10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends: []BackendConfig{
					{URL: "http://localhost:8081"},
				},
			},
			wantErr:     true,
			errContains: "health check interval must be positive",
		},
		{
			name: "no backends",
			config: &Config{
				Port:                8080,
				HealthCheckInterval: 10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends:            []BackendConfig{},
			},
			wantErr:     true,
			errContains: "at least one backend is required",
		},
		{
			name: "invalid backend URL",
			config: &Config{
				Port:                8080,
				HealthCheckInterval: 10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends: []BackendConfig{
					{URL: ""},
				},
			},
			wantErr:     true,
			errContains: "URL is required",
		},
		{
			name: "invalid backend weight",
			config: &Config{
				Port:                8080,
				HealthCheckInterval: 10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends: []BackendConfig{
					{
						URL:    "http://localhost:8081",
						Weight: -1,
					},
				},
			},
			wantErr:     true,
			errContains: "weight cannot be negative",
		},
		{
			name: "invalid backend maxConns",
			config: &Config{
				Port:                8080,
				HealthCheckInterval: 10 * time.Second,
				Strategy:            serverpool.RoundRobin,
				Backends: []BackendConfig{
					{
						URL:      "http://localhost:8081",
						MaxConns: -1,
					},
				},
			},
			wantErr:     true,
			errContains: "maxConns cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		configYAML  string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid configuration",
			configYAML: `
port: 8080
healthCheckInterval: 10s
strategy: 0
backends:
  - url: "http://localhost:8081"
    weight: 1
    maxConns: 100
  - url: "http://localhost:8082"
    weight: 2
    maxConns: 200
`,
			wantErr: false,
		},
		{
			name: "invalid port",
			configYAML: `
port: -1
healthCheckInterval: 10s
strategy: 0
backends:
  - url: "http://localhost:8081"
`,
			wantErr:     true,
			errContains: "invalid port number",
		},
		{
			name: "invalid health check interval",
			configYAML: `
port: 8080
healthCheckInterval: -10s
strategy: 0
backends:
  - url: "http://localhost:8081"
`,
			wantErr:     true,
			errContains: "health check interval must be positive",
		},
		{
			name: "no backends",
			configYAML: `
port: 8080
healthCheckInterval: 10s
strategy: 0
backends: []
`,
			wantErr:     true,
			errContains: "at least one backend is required",
		},
		{
			name: "invalid backend URL",
			configYAML: `
port: 8080
healthCheckInterval: 10s
strategy: 0
backends:
  - url: ""
`,
			wantErr:     true,
			errContains: "URL is required",
		},
		{
			name: "invalid backend weight",
			configYAML: `
port: 8080
healthCheckInterval: 10s
strategy: 0
backends:
  - url: "http://localhost:8081"
    weight: -1
`,
			wantErr:     true,
			errContains: "weight cannot be negative",
		},
		{
			name: "invalid backend maxConns",
			configYAML: `
port: 8080
healthCheckInterval: 10s
strategy: 0
backends:
  - url: "http://localhost:8081"
    maxConns: -1
`,
			wantErr:     true,
			errContains: "maxConns cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config file
			configPath := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configPath, []byte(tt.configYAML), 0644)
			assert.NoError(t, err)

			// Try to load the configuration
			cfg, err := LoadFromFile(configPath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, cfg)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, cfg)

			// Verify the configuration
			assert.Equal(t, 8080, cfg.Port)
			assert.Equal(t, 10*time.Second, cfg.HealthCheckInterval)
			assert.Equal(t, serverpool.RoundRobin, cfg.Strategy)
			assert.Len(t, cfg.Backends, 2)
		})
	}
}

func TestLoadFromFileWithInvalidFile(t *testing.T) {
	// Test with non-existent file
	_, err := LoadFromFile("nonexistent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")

	// Test with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err = os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644)
	assert.NoError(t, err)

	_, err = LoadFromFile(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

func TestLoadFromEnv(t *testing.T) {
	// Test that LoadFromEnv returns an error when no backends are provided
	cfg, err := LoadFromEnv()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one backend is required")
	assert.Nil(t, cfg)
}
