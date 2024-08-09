# Stage 1: Build
FROM golang:latest AS builder
LABEL authors="manuelantoniorojasramos"

# Set the Current Working Directory inside the container
WORKDIR /home/builder/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install certificates to allow HTTPS connections
RUN apk --no-cache add ca-certificates


# Set the Current Working Directory inside the container
WORKDIR /usr/local/bin

# Copy the binary from the builder stage
COPY --from=builder /home/builder/app/server /usr/local/bin/server
COPY templates ./templates

# Expose port (change if your application uses a different port)
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/usr/local/bin/server"]
