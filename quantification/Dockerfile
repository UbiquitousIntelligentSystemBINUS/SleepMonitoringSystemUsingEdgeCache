# Start with the Go alpine image
FROM golang:alpine AS builder

ENV REDIS_ADDR="localhost:6379"
ENV DB_HOST="localhost"
ENV DB_PORT=5432
ENV DB_USER="postgres" \
    DB_PASSWORD="postgres" \
    DB_NAME="sleep_monitoring" \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
# COPY go.mod go.sum ./

COPY . .
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 9001 to the outside world
EXPOSE 9001

# Command to run the executable
CMD ["./main"]
