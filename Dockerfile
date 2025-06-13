# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3

# Add necessary tools
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/main .

# Copy .env file (will be overridden by docker-compose env vars if provided)
COPY .env.example .env

# Change ownership of the working directory
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
