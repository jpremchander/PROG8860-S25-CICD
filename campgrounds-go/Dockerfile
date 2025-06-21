FROM golang:1.21-alpine AS builder

# Install ca-certificates and tzdata
RUN apk add --no-cache ca-certificates tzdata git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Set Go proxy and sum database
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o yelpcamp-app .

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Set working directory
WORKDIR /app

# Create uploads directory with proper permissions
RUN mkdir -p static/uploads && chmod 755 static/uploads

# Copy binary from builder stage
COPY --from=builder /app/yelpcamp-app .

# Copy static files and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Copy environment template (JWT will be auto-generated)
COPY --from=builder /app/.env.template ./.env

# Set environment for production
ENV GIN_MODE=release
ENV ENVIRONMENT=production

# Set timezone
RUN cp /usr/share/zoneinfo/UTC /etc/localtime

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Change ownership of app directory and uploads
RUN chown -R appuser:appgroup /app && \
    chmod -R 755 /app/static

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Run the application (JWT will be auto-generated on startup)
CMD ["./yelpcamp-app"]
