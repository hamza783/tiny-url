# run this command from the tiny-url directory. This is needed because `shared` is not its own github repo.
# to create an image: docker build -t url-shortening-service -f url-shortening-service/Dockerfile .
# to run the container: docker run -p 8081:8081 url-shortening-service 

# Go version
FROM golang:1.23-alpine AS builder

# Install git and CA certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy gp.mod and go.sum for dependency caching
COPY url-shortening-service/go.mod url-shortening-service/go.sum ./url-shortening-service/
# Copy go.mod for shared module
COPY shared/go.mod ./shared/go.mod

WORKDIR /app/url-shortening-service
RUN go mod download

# Copy source code
WORKDIR /app
# Only copy url-shortening-service and shared because only these are needed for this build
COPY url-shortening-service ./url-shortening-service
COPY shared ./shared

# Build
WORKDIR /app/url-shortening-service
RUN go build -o /url-shortening-service ./cmd/main.go

# Run
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /url-shortening-service .

# Expose the port your service listens on
EXPOSE 8081

# Run the binary
CMD ["./url-shortening-service"]