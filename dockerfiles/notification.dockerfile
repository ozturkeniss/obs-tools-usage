# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the notification service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/notification-service cmd/notification/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/notification-service .

# Create non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

# Expose port
EXPOSE 8084

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8084/health || exit 1

# Run the application
CMD ["./notification-service"]
