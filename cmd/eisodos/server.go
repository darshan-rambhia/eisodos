package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darshan-rambhia/eisodos/internal/backend"
)

// LoadBalancer interface defines the methods required by the server
type LoadBalancer interface {
	GetPort() int
	Start() error
	Stop(ctx context.Context) error
	GetServerPoolSize() int
	GetBackends() []backend.Backend
}

// Server represents the load balancer server
type Server struct {
	lb LoadBalancer
}

// NewServer creates a new server instance
func NewServer(lb LoadBalancer) *Server {
	return &Server{lb: lb}
}

// Start starts the server and handles shutdown
func (s *Server) Start() error {
	// Create a channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Starting load balancer on port %d", s.lb.GetPort())
		serverErrors <- s.lb.Start()
	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a server error or a shutdown signal
	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		log.Printf("Received signal %v, initiating shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.lb.Stop(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
			return err
		}
	}

	return nil
}
