# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Create database directory
RUN mkdir -p database

# Set environment variables
ENV CGO_ENABLED=1

# Run the application
CMD ["./main"]
