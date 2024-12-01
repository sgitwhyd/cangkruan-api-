# Stage 1: Build the Go binary
FROM golang:1.22 AS builder

WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the application
COPY . ./

# Build the application with necessary environment variables for a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Debug: List files in builder stage
RUN ls -l /app

# Stage 2: Create a minimal runtime image
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies (if needed)
RUN apk --no-cache add ca-certificates

# Copy the binary and configs from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/internal/configs ./internal/configs

# Ensure the binary has executable permissions
RUN chmod +x /app/main

# Expose the application port
EXPOSE 9876

# Run the binary
CMD ["./main"]
