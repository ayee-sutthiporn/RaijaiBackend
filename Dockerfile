# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o main ./cmd/api/main.go

# Stage 2: Run
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port (default 8080)
EXPOSE 8088

# Run the binary
CMD ["./main"]
