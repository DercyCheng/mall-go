FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./services/user-service/cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Shanghai

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .

# Copy configuration
COPY --from=builder /app/services/user-service/configs/config.yaml ./configs/

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./user-service"]