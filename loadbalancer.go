package eisodos

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/darshan-rambhia/eisodos/config"
	"github.com/darshan-rambhia/eisodos/internal/backend"
	"github.com/darshan-rambhia/eisodos/internal/serverpool"
)

// LoadBalancer represents the main load balancer instance
type LoadBalancer struct {
	serverPool serverpool.ServerPool
	server     *http.Server
	mu         sync.RWMutex
}

// LoadBalancerBuilder provides a fluent interface for building a LoadBalancer
type LoadBalancerBuilder struct {
	config     *config.Config
	backends   []backend.Backend
	serverPool serverpool.ServerPool
}

// NewLoadBalancerBuilder creates a new LoadBalancerBuilder
func NewLoadBalancerBuilder() *LoadBalancerBuilder {
	return &LoadBalancerBuilder{
		config: config.DefaultConfig(),
	}
}

// WithConfig sets the configuration for the load balancer
func (b *LoadBalancerBuilder) WithConfig(cfg *config.Config) *LoadBalancerBuilder {
	b.config = cfg
	return b
}

// WithPort sets the port for the load balancer
func (b *LoadBalancerBuilder) WithPort(port int) *LoadBalancerBuilder {
	b.config.Port = port
	return b
}

// WithHealthCheckInterval sets the health check interval
func (b *LoadBalancerBuilder) WithHealthCheckInterval(interval time.Duration) *LoadBalancerBuilder {
	b.config.HealthCheckInterval = interval
	return b
}

// WithStrategy sets the load balancing strategy
func (b *LoadBalancerBuilder) WithStrategy(strategy serverpool.LBStrategy) *LoadBalancerBuilder {
	b.config.Strategy = strategy
	return b
}

// WithBackend adds a backend to the load balancer
func (b *LoadBalancerBuilder) WithBackend(url *url.URL, proxy *httputil.ReverseProxy) *LoadBalancerBuilder {
	b.backends = append(b.backends, backend.NewBackend(url, proxy))
	return b
}

// Build creates and returns a new LoadBalancer instance
func (b *LoadBalancerBuilder) Build() (*LoadBalancer, error) {
	pool, err := serverpool.NewServerPool(b.config.Strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create server pool: %w", err)
	}

	lb := &LoadBalancer{
		serverPool: pool,
	}

	// Create HTTP server
	lb.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", b.config.Port),
		Handler: lb,
	}

	// Add backends
	for _, b := range b.backends {
		lb.AddBackend(b)
	}

	// Start health check routine
	go lb.startHealthCheck(b.config.HealthCheckInterval)

	return lb, nil
}

// Config holds the configuration for the load balancer
type Config struct {
	Port                int
	HealthCheckInterval time.Duration
	Strategy            serverpool.LBStrategy
}

// NewLoadBalancer creates a new load balancer instance
func NewLoadBalancer(cfg Config) (*LoadBalancer, error) {
	pool, err := serverpool.NewServerPool(cfg.Strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create server pool: %w", err)
	}

	lb := &LoadBalancer{
		serverPool: pool,
	}

	// Create HTTP server
	lb.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: lb,
	}

	// Start health check routine
	go lb.startHealthCheck(cfg.HealthCheckInterval)

	return lb, nil
}

// ServeHTTP implements the http.Handler interface
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	peer := lb.serverPool.GetNextValidPeer()
	if peer == nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer.Serve(w, r)
}

// startHealthCheck runs the health check routine
func (lb *LoadBalancer) startHealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		serverpool.HealthCheck(ctx, lb.serverPool)
		cancel()
	}
}

// Start starts the load balancer server
func (lb *LoadBalancer) Start() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	return lb.server.ListenAndServe()
}

// Stop gracefully shuts down the load balancer
func (lb *LoadBalancer) Stop(ctx context.Context) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	return lb.server.Shutdown(ctx)
}

// AddBackend adds a new backend to the server pool
func (lb *LoadBalancer) AddBackend(b backend.Backend) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.serverPool.AddBackend(b)
}

// GetBackends returns all backends in the server pool
func (lb *LoadBalancer) GetBackends() []backend.Backend {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return lb.serverPool.GetBackends()
}

// GetServerPoolSize returns the current size of the server pool
func (lb *LoadBalancer) GetServerPoolSize() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return lb.serverPool.GetServerPoolSize()
}

// GetPort returns the port the load balancer is listening on
func (lb *LoadBalancer) GetPort() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	addr := lb.server.Addr
	if addr == "" {
		return 0
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0
	}
	return p
}
