# run this command from the tiny-url directory. This is needed because `shared` is not its own github repo.
# to create an image: docker build -t url-redirection-service -f url-redirection-service/Dockerfile .
# to run the container: docker run -p 9000:9000 url-redirection-service 

# Go version
FROM golang:1.23-alpine AS builder

# Install git and CA certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy gp.mod and go.sum for dependency caching
COPY url-redirection-service/go.mod url-redirection-service/go.sum ./url-redirection-service/
# Copy go.mod for shared module
COPY shared/go.mod ./shared/go.mod

WORKDIR /app/url-redirection-service
RUN go mod download

# Copy source code
WORKDIR /app
# Only copy url-redirection-service and shared because only these are needed for this build
COPY url-redirection-service ./url-redirection-service
COPY shared ./shared

# Build
WORKDIR /app/url-redirection-service
RUN go build -o /url-redirection-service ./cmd/main.go

# Run
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /url-redirection-service .

# Expose the port your service listens on
EXPOSE 9000

# Run the binary
CMD ["./url-redirection-service"]