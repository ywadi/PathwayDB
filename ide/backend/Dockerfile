# Stage 1: Builder
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY ide/backend/go.mod ide/backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY ide/backend/ .

# Run backend tests (if any exist)
RUN if [ -d "./tests" ] || find . -name "*_test.go" | grep -q .; then go test ./... -v; fi

# Build the backend-server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend-server .

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

# Copy the static binary from the builder stage
COPY --from=builder /backend-server .

# Create a non-root user
RUN addgroup -S appgroup && adduser -S -G appgroup appuser
RUN chown -R appuser:appgroup /app
USER appuser

# Expose the port
EXPOSE 8081

# Set the entrypoint
CMD ["./backend-server"]
