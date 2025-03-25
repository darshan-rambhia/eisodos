package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/darshan-rambhia/eisodos"
	"github.com/darshan-rambhia/eisodos/internal/serverpool"
)

// LoadFromCLI creates a load balancer from command line arguments
func LoadFromCLI(port int, healthCheckInterval time.Duration, strategy string, backendURLs []string) (LoadBalancer, error) {
	if len(backendURLs) == 0 {
		return nil, fmt.Errorf("please provide at least one backend URL")
	}

	builder := eisodos.NewLoadBalancerBuilder().
		WithPort(port).
		WithHealthCheckInterval(healthCheckInterval).
		WithStrategy(serverpool.ParseStrategy(strategy))

	for _, backendURL := range backendURLs {
		url, err := url.Parse(backendURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse backend URL %s: %w", backendURL, err)
		}

		if url.Scheme == "" || url.Host == "" {
			return nil, fmt.Errorf("failed to parse backend URL %s: missing scheme or host", backendURL)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			http.Error(w, "Proxy error", http.StatusBadGateway)
		}

		builder.WithBackend(url, proxy)
	}

	return builder.Build()
}
