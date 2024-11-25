# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache \
    tzdata \
    sqlite \
    ca-certificates

# Set the timezone
ENV TZ=Asia/Jakarta

# Create app directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy any additional required files (like .env.example)
COPY .env.example .env

# Create a directory for SQLite database
RUN mkdir -p /app/database && \
    chown -R nobody:nobody /app

# Use non-root user
USER nobody

# Command to run the application
CMD ["./main"]
