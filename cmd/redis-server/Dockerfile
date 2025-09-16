# Stage 1: Builder
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the redis-server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /redis-server ./cmd/redis-server/main.go

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

# Copy the static binary from the builder stage
COPY --from=builder /redis-server .

# Create a non-root user and data directory
RUN addgroup -S appgroup && adduser -S -G appgroup appuser
RUN mkdir -p /app/data && chown -R appuser:appgroup /app
USER appuser

# Expose the port
EXPOSE 6379

# Set the entrypoint
CMD ["./redis-server"]
