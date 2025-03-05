# Use official Go image for building, then switch to Alpine for smaller final image
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for Go module dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sharebin main.go

# Use a minimal Alpine image for the runtime
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/sharebin /app/

# Create directories for uploads, data, templates, and static files
RUN mkdir -p /app/uploads /app/data /app/templates /app/static

# Set permissions for the directories
RUN chmod -R 755 /app/uploads /app/data /app/templates /app/static

# Expose port 80 (default HTTP port)
EXPOSE 80

# Set environment variables (optional, for configuration)
ENV SHAREBIN_PORT=80

# Run the application
CMD ["./sharebin"]
