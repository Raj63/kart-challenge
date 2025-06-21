FROM --platform=linux/amd64 golang:1.23-alpine AS builder

WORKDIR /app/backend-challenge/services/coupons

# Copy go mod and sum
COPY backend-challenge/services/coupons/go.mod .
COPY backend-challenge/services/coupons/go.sum .

# Copy shared library using full relative path
COPY backend-challenge/library /app/backend-challenge/library

RUN go mod download

# Copy the rest of the code
COPY backend-challenge/services/coupons/. .

# Download dependencies
RUN go mod download

# Build the processor binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./coupons-processor ./cmd/processor/main.go

# Final runtime image
FROM --platform=linux/amd64 alpine:latest

# Install ca-certificates for MongoDB TLS
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/backend-challenge/services/coupons/coupons-processor ./
COPY --from=builder /app/backend-challenge/services/coupons/config.json ./
# COPY --from=builder /app/backend-challenge/services/coupons/data /data
RUN chmod +x /coupons-processor
RUN mkdir -p /logs

VOLUME ["/app/data"]

ENTRYPOINT ["./coupons-processor"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 CMD pgrep coupons-processor || exit 1 