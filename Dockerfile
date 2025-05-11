# ---- Build Stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git (required for go mod)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- Run Stage ----
FROM alpine:latest
WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/server ./server

# Copy migrations, config, and any static files if needed
COPY --from=builder /app/db ./db
COPY --from=builder /app/auth_model.conf ./auth_model.conf
COPY --from=builder /app/.env ./

# Expose port (default 8000)
EXPOSE 8000

# Run the app
CMD ["./server"]
