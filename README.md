# Advanced Load Balancer Implementation Project

A comprehensive project to demonstrate and strengthen advanced Go programming, Kubernetes, and distributed systems concepts through the implementation of a load balancer. This project serves as both a learning journey and a portfolio piece.

## Project Overview

This project implements a production-grade load balancer with various algorithms and deployment strategies, showcasing advanced Go programming patterns, Kubernetes expertise, and distributed systems concepts.

## Learning Objectives & Implementation Checklist

### Go Programming Excellence

- [x] Advanced Concurrency Patterns
  - [x] Implementation of worker pools
  - [x] Channel-based communication patterns
  - [x] Context usage for cancellation and timeouts
  - [x] Goroutine lifecycle management
  - [ ] Rate limiting implementation
  
- [x] Clean Code Architecture
  - [x] Hexagonal/Clean Architecture implementation
  - [x] Dependency injection patterns
  - [x] Interface segregation
  - [x] Error handling best practices
  - [x] Middleware implementation

- [x] Testing Excellence
  - [x] Unit testing with table-driven tests
  - [x] Integration testing
  - [ ] Benchmark testing
  - [ ] Fuzzing tests
  - [x] Mock implementations
  
- [x] Performance Optimization
  - [x] Memory management and optimization
  - [ ] CPU profiling and optimization
  - [x] Connection pooling
  - [ ] Caching strategies

### Data Structures & Algorithms

- [x] Load Balancing Algorithms Implementation
  - [x] Round Robin
  - [ ] Weighted Round Robin
  - [x] Least Connections
  - [ ] Consistent Hashing
  - [ ] IP Hash-based routing

- [ ] Custom Data Structures
  - [ ] Thread-safe priority queue
  - [ ] Concurrent hash map
  - [ ] Custom heap implementation
  - [ ] LRU cache implementation

### Kubernetes & Cloud Native

- [ ] Gateway API Implementation
  - [ ] HTTPRoute implementation
  - [ ] TCPRoute implementation
  - [ ] TLSRoute implementation
  - [ ] Custom policy attachments

- [ ] Kubernetes Integration
  - [ ] Custom Resource Definitions (CRDs)
  - [ ] Custom controller implementation
  - [ ] Operator pattern
  - [ ] Service discovery integration
  
- [x] Observability
  - [ ] Prometheus metrics integration
  - [ ] OpenTelemetry tracing
  - [x] Structured logging
  - [x] Health check endpoints

### Advanced Features

- [ ] Security Implementation
  - [ ] TLS termination
  - [ ] Certificate management
  - [ ] Rate limiting
  - [ ] WAF-like features

- [x] High Availability
  - [x] Leader election
  - [x] State synchronization
  - [x] Failover mechanisms
  - [x] Session persistence

- [x] Traffic Management
  - [x] Circuit breaking
  - [x] Retry mechanisms
  - [x] Timeout handling
  - [ ] Traffic splitting

## Project Structure

```text
├── cmd/
│   └── eisodos/           # Main application entry point
├── internal/
│   ├── backend/          # Backend server implementation
│   ├── serverpool/       # Load balancing strategies
│   └── config/           # Configuration management
├── buildscripts/
│   └── scripts/          # Build and test scripts
└── test/                 # Test files
```

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/darshan-rambhia/eisodos.git
   cd eisodos
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go run buildscripts/scripts/test.go
   ```

4. Build the application:
   ```bash
   go build -o eisodos cmd/eisodos/main.go
   ```

5. Run the load balancer:
   ```bash
   ./eisodos --port 8080 --backends http://localhost:8081,http://localhost:8082
   ```

## Development Workflow

1. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and run tests:
   ```bash
   go run buildscripts/scripts/test.go
   ```

3. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

4. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```

5. Create a pull request on GitHub.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

MIT License

## Author

Darshan Rambhia
