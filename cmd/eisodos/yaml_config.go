package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/darshan-rambhia/eisodos"
	"github.com/darshan-rambhia/eisodos/config"
)

// LoadFromYAML creates a load balancer from a YAML configuration file
func LoadFromYAML(path string) (LoadBalancer, error) {
	cfg, err := config.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	builder := eisodos.NewLoadBalancerBuilder().
		WithConfig(cfg)

	for _, backend := range cfg.Backends {
		url, err := url.Parse(backend.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse backend URL %s: %w", backend.URL, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "Proxy error", http.StatusBadGateway)
		}

		builder.WithBackend(url, proxy)
	}

	return builder.Build()
}
