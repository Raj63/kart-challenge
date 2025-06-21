# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod and sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the processor binary
RUN go build -o coupons-processor ./cmd/processor/main.go

# Final runtime image
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for MongoDB TLS
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/coupons-processor ./
COPY config.json ./

# Create data directory for mounting
RUN mkdir -p /app/data/add /app/data/remove

VOLUME ["/app/data"]

ENTRYPOINT ["./coupons-processor"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 CMD pgrep coupons-processor || exit 1 