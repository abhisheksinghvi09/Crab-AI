# Multi-stage build for optimized production image
# Build stage
FROM golang:1.23-alpine AS builder

# Add necessary build tools and security updates
RUN apk add --no-cache git ca-certificates tzdata && \
    apk upgrade

# Create non-root user for build
RUN adduser -D -g '' builduser

# Set working directory
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with security flags and optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath \
    -ldflags="-w -s -X main.version=$(date +%Y%m%d-%H%M%S)" \
    -a -installsuffix cgo \
    -o /app/main cmd/main.go

# Final stage - minimal production image
FROM alpine:3.19

# Install runtime dependencies and security updates
RUN apk --no-cache add ca-certificates tzdata wget && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G appgroup appuser

# Set working directory
WORKDIR /app

# Copy binary from builder with correct ownership
COPY --from=builder --chown=appuser:appgroup /app/main .

# Copy configuration files with correct ownership
COPY --from=builder --chown=appuser:appgroup /app/config ./config

# Copy email templates with correct ownership
COPY --from=builder --chown=appuser:appgroup /app/internal/email_templates ./internal/email_templates

# Create necessary directories with correct permissions
RUN mkdir -p /app/logs /app/tmp && \
    chown -R appuser:appgroup /app && \
    chmod -R 755 /app && \
    chmod +x /app/main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Set environment variables for production
ENV GIN_MODE=release
ENV PORT=3000

# Add labels for better container management
LABEL maintainer="Crab-AI Team" \
      description="Crab-AI Backend Service" \
      version="1.0.0"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run the application
ENTRYPOINT ["./main"]