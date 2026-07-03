# --- STAGE 1: Build the Go Binary ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Compile the Go application into a single static binary named 'server'
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# --- STAGE 2: Run the Compiled Binary ---
FROM alpine:latest

RUN apk add ffmpeg ca-certificates

WORKDIR /app

# Copy only the compiled execution binary from Stage 1
COPY --from=builder /app/server .

EXPOSE 8080

# Run the binary
CMD ["./server"]