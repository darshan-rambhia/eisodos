package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	// Command line flags
	configFile := flag.String("config", "", "Path to YAML configuration file")
	port := flag.Int("port", 8080, "Port to listen on")
	healthCheckInterval := flag.Duration("health-check-interval", 10*time.Second, "Health check interval")
	strategy := flag.String("strategy", "round-robin", "Load balancing strategy (round-robin or least-connected)")
	flag.Parse()

	var lb LoadBalancer
	var err error

	if *configFile != "" {
		// Load from YAML config file
		lb, err = LoadFromYAML(*configFile)
	} else {
		// Use command line arguments
		backendURLs := flag.Args()
		lb, err = LoadFromCLI(*port, *healthCheckInterval, *strategy, backendURLs)
	}

	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}

	// Create and start the server
	server := NewServer(lb)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
