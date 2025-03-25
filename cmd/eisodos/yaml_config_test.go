package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromYAML(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config file
			configPath := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configPath, []byte(tt.configYAML), 0644)
			assert.NoError(t, err)

			// Try to load the configuration
			lb, err := LoadFromYAML(configPath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, lb)

			// Verify the load balancer configuration
			assert.Equal(t, 8080, lb.GetPort())
			assert.Equal(t, 2, lb.GetServerPoolSize())
		})
	}
}

func TestLoadFromYAMLWithInvalidFile(t *testing.T) {
	// Test with non-existent file
	_, err := LoadFromYAML("nonexistent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")

	// Test with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err = os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644)
	assert.NoError(t, err)

	_, err = LoadFromYAML(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}
