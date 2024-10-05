# Use an official Go image as the base image
FROM golang:1.23-alpine

# Set Go environment variables
ENV GO111MODULE=on

# Copy go.mod and go.sum to the root of the container
COPY go.mod go.sum ./

# Download Go modules and dependencies
RUN go mod download

# Copy the entire project to the root of the container
COPY . .

# Build the Go application using the existing structure
RUN go build -o app ./main.go

# Set the entry point command to run the Go app
CMD ["./app"]

