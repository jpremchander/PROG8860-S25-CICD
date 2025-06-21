# Simplified Dockerfile that avoids Git issues
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy everything first
COPY . .

# Set Go proxy to avoid Git authentication issues
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Download dependencies with retry
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o yelpcamp-app .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary and assets
COPY --from=builder /app/yelpcamp-app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/api/health || exit 1

CMD ["./yelpcamp-app"]