# Build stage
FROM golang:1.22.4-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server.go

# Runtime stage
FROM alpine:latest

# Install SQLite runtime
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port (Cloud Run will set PORT env var)
EXPOSE 8080

# Run the server
CMD ["./server"]

