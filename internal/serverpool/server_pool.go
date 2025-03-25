package serverpool

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/darshan-rambhia/eisodos/internal/backend"
)

type ServerPool interface {
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	AddBackend(backend.Backend)
	GetServerPoolSize() int
}

func HealthCheck(ctx context.Context, s ServerPool) {
	aliveChannel := make(chan bool, 1)

	for _, b := range s.GetBackends() {
		b := b
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		status := "up"
		go backend.IsBackendAlive(requestCtx, aliveChannel, b.GetURL())

		select {
		case <-ctx.Done():
			slog.Info("Gracefully shutting down health check")
			return
		case alive := <-aliveChannel:
			b.SetAlive(alive)
			if !alive {
				status = "down"
			}
		}
		slog.Debug(
			"URL Status",
			"URL", b.GetURL().String(),
			"status", status,
		)
	}
}

func NewServerPool(strategy LBStrategy) (ServerPool, error) {
	switch strategy {
	case RoundRobin:
		return &roundRobinServerPool{
			backends: make([]backend.Backend, 0),
			current:  0,
		}, nil
	case LeastConnected:
		return &lcServerPool{
			backends: make([]backend.Backend, 0),
		}, nil
	default:
		return nil, errors.New("invalid strategy")
	}
}
