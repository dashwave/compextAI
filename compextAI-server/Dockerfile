# Use an official Go image as the builder
FROM golang:1.22.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies (this layer is cached unless go.mod or go.sum changes)
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main *.go

# Use a minimal base image for the final executable
FROM alpine:latest

WORKDIR /root

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8888

# Command to run the executable
CMD ["./main"]