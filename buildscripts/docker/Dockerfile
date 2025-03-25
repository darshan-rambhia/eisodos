FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Create a non-root user
RUN adduser -D -g '' appuser

# Set up Go module cache
ENV GOPATH /go
ENV GOMODCACHE /go/pkg/mod

# Create and set the working directory
WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/.out/bin/eisodos ./cmd/eisodos

# Create a minimal production image
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/.out/bin/eisodos .

# Use non-root user
USER appuser

EXPOSE 8080
CMD ["./eisodos"] 