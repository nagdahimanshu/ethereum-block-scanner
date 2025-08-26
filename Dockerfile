FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary from cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o scanner ./cmd/block-scanner/main.go

FROM alpine:latest

# Install CA certificates for HTTPS (needed by Ethereum RPC and Kafka)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/scanner .

# Set the entrypoint to your scanner binary
ENTRYPOINT ["./scanner"]
