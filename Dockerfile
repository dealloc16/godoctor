# Use the official Golang image to create a build artifact.
# This is the "builder" stage.
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o bin/godoctor cmd/godoctor/main.go


# Start from a new, minimal image to reduce the image size.
# This is the "final" stage.
FROM golang:1.24-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/godoctor .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./godoctor"]
