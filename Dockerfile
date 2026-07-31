# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy migrations directory
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

# Default env (can be overridden via docker-compose or env vars)
ENV APP_ENV=production
ENV APP_PORT=8080
ENV APP_HOST=0.0.0.0
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_NAME=tms
ENV DB_USER=tms_user
ENV DB_PASSWORD=secret
ENV DB_SSLMODE=disable
ENV JWT_SECRET=change-this-in-production

ENTRYPOINT ["./server"]
