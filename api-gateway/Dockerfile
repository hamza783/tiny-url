# run this command from the tiny-url directory. This is needed because `shared` is not its own github repo.
# to create an image: docker build -t api-gateway -f api-gateway/Dockerfile .
# to run the container: docker run -p 8080:8080 api-gateway  

# Go version
FROM golang:1.23-alpine AS builder

# Install git and CA certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy gp.mod and go.sum for dependency caching
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
# Copy go.mod for shared module
COPY shared/go.mod ./shared/go.mod

WORKDIR /app/api-gateway
RUN go mod download

# Copy source code
WORKDIR /app
# Only copy api-gateway and shared because only these are needed for this build
COPY api-gateway ./api-gateway
COPY shared ./shared

# Build
WORKDIR /app/api-gateway
RUN go build -o /api-gateway ./cmd/main.go

# Run
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /api-gateway .

# Expose the port your service listens on
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]