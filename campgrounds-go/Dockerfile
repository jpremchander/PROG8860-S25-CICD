# Base stage with common dependencies
FROM golang:1.23-alpine AS base

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates curl

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Development stage
FROM base AS development

# Install air for hot reloading (development only)
RUN go install github.com/cosmtrek/air@latest

# Expose development port
EXPOSE 3001

# Set development environment
ENV GIN_MODE=debug
ENV PORT=3001

# Create non-root user for development
RUN addgroup -g 1001 -S devgroup && \
    adduser -u 1001 -S devuser -G devgroup

# Change ownership
RUN chown -R devuser:devgroup /app

USER devuser

# Development health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3001/health || exit 1

# Run with hot reload in development
CMD ["air", "-c", ".air.toml"]

# Production build stage
FROM base AS builder

# Build the application for production
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o campgrounds-go main.go

# Production stage
FROM alpine:latest AS production

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user for production
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/campgrounds-go .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose production port
EXPOSE 3000

# Set production environment
ENV GIN_MODE=release
ENV PORT=3000

# Production health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Run the application
CMD ["./campgrounds-go"]
