package main

import (
	"context"
	"testing"
	"time"

	"github.com/darshan-rambhia/eisodos/internal/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockLoadBalancer is a mock implementation of the LoadBalancer interface
type mockLoadBalancer struct {
	mock.Mock
}

func (m *mockLoadBalancer) GetPort() int {
	args := m.Called()
	return args.Int(0)
}

func (m *mockLoadBalancer) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockLoadBalancer) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockLoadBalancer) GetServerPoolSize() int {
	args := m.Called()
	return args.Int(0)
}

func (m *mockLoadBalancer) GetBackends() []backend.Backend {
	args := m.Called()
	return args.Get(0).([]backend.Backend)
}

func TestServerStart(t *testing.T) {
	// Create a mock load balancer
	mockLB := new(mockLoadBalancer)
	mockLB.On("GetPort").Return(8080)
	mockLB.On("Start").Return(nil)
	mockLB.On("Stop", mock.Anything).Return(nil)
	mockLB.On("GetServerPoolSize").Return(1)
	mockLB.On("GetBackends").Return([]backend.Backend{})

	// Create a server with the mock load balancer
	server := NewServer(mockLB)

	// Start the server in a goroutine
	go func() {
		err := server.Start()
		assert.NoError(t, err)
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Verify that the load balancer was started
	mockLB.AssertCalled(t, "GetPort")
	mockLB.AssertCalled(t, "Start")
}

func TestServerStartWithError(t *testing.T) {
	// Create a mock load balancer that returns an error on start
	mockLB := new(mockLoadBalancer)
	mockLB.On("GetPort").Return(8080)
	mockLB.On("Start").Return(assert.AnError)
	mockLB.On("GetServerPoolSize").Return(1)
	mockLB.On("GetBackends").Return([]backend.Backend{})

	// Create a server with the mock load balancer
	server := NewServer(mockLB)

	// Start the server and expect an error
	err := server.Start()
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
