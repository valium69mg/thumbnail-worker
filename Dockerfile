# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o thumbnail-worker main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Copy binary
COPY --from=builder /app/thumbnail-worker .

RUN apt-get update \
    && apt-get install -y ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*


# Run worker
CMD ["./thumbnail-worker"]
