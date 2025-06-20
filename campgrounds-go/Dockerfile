# Development Dockerfile - Fixed for missing go.sum
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first
COPY go.mod ./
# Copy go.sum if it exists, otherwise create empty file
COPY go.su[m] ./
RUN if [ ! -f go.sum ]; then touch go.sum; fi

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Ensure go.sum is properly generated
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o yelpcamp-app .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/yelpcamp-app .

# Copy static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/public ./public

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Change ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run the binary
CMD ["./yelpcamp-app"]
