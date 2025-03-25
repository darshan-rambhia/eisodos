package config

import (
	"fmt"
	"os"
	"time"

	"github.com/darshan-rambhia/eisodos/internal/serverpool"
	"gopkg.in/yaml.v3"
)

// Config represents the load balancer configuration
type Config struct {
	Port                int                   `yaml:"port"`
	HealthCheckInterval time.Duration         `yaml:"healthCheckInterval"`
	Strategy            serverpool.LBStrategy `yaml:"strategy"`
	Backends            []BackendConfig       `yaml:"backends"`
}

// BackendConfig represents a backend server configuration
type BackendConfig struct {
	URL      string `yaml:"url"`
	Weight   int    `yaml:"weight,omitempty"`
	MaxConns int    `yaml:"maxConns,omitempty"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:                8080,
		HealthCheckInterval: 10 * time.Second,
		Strategy:            serverpool.RoundRobin,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Port)
	}

	if c.HealthCheckInterval <= 0 {
		return fmt.Errorf("health check interval must be positive: %v", c.HealthCheckInterval)
	}

	if len(c.Backends) == 0 {
		return fmt.Errorf("at least one backend is required")
	}

	for i, backend := range c.Backends {
		if backend.URL == "" {
			return fmt.Errorf("backend %d: URL is required", i)
		}
		if backend.Weight < 0 {
			return fmt.Errorf("backend %d: weight cannot be negative", i)
		}
		if backend.MaxConns < 0 {
			return fmt.Errorf("backend %d: maxConns cannot be negative", i)
		}
	}

	return nil
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	config := DefaultConfig()

	// TODO: Implement environment variable loading
	// This would read from environment variables like:
	// EISODOS_PORT
	// EISODOS_HEALTH_CHECK_INTERVAL
	// EISODOS_STRATEGY
	// EISODOS_BACKENDS (comma-separated list of URLs)

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}
