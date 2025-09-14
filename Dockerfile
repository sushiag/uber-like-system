# Start from the official Go image
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o server ./cmd/server

# --- runtime stage ---
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]