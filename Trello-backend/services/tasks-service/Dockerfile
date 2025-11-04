# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder  

# Set the working directory for the build
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code to the working directory
COPY src/ ./src

# Build the Go application, outputting the binary named tasks-service
RUN GOOS=linux GOARCH=amd64 go build -o tasks-service ./src/main.go

# Stage 2: Create a smaller image for the runtime
FROM alpine:latest  

# Set the working directory for the runtime
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/tasks-service .

# Set the permissions for the binary
RUN chmod +x tasks-service

# Command to run the binary when the container starts
CMD ["./tasks-service"]
