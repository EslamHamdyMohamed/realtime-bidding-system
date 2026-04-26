# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum if they exist
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
WORKDIR /app/services/api-gateway
RUN go build -o /api-gateway ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /api-gateway /api-gateway
COPY --from=builder /app/services/api-gateway/.env.example /.env

EXPOSE 8080

CMD ["/api-gateway"]
