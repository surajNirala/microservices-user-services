# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS build

# Set environment variables for Go
ENV GO111MODULE=on

# Create a directory for the app
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go app
RUN go build -o user_services .

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set a working directory in the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=build /app/user_services .

# Expose the necessary port (if the service runs on 8080)
EXPOSE 9091

# Run the application
CMD ["./user_services"]
