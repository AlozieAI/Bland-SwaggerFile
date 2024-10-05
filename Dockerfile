# Use an official Go image as the base image
FROM golang:1.23-alpine

# Install git to ensure go mod can download dependencies
RUN apk add --no-cache git

# Set Go environment variables
ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOPATH=/go
ENV GOMODCACHE=/go/pkg/mod

# Copy go.mod and go.sum to the root of the container
COPY go.mod go.sum ./

# Download Go modules and dependencies
RUN go mod download

# Copy the entire project to the root of the container
COPY . .

# Build the Go application
RUN go build -o app ./main.go

# Set the entry point command to run the Go app
CMD ["./app"]
