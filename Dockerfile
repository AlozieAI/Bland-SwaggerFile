# Use an official Go image as the base image
FROM golang:1.23-alpine

# Install git to ensure go mod can download dependencies
RUN apk add --no-cache git

# Set Go environment variables
ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOPATH=/go
ENV GOMODCACHE=/go/pkg/mod

# Create directory structure and set the working directory to /go/src/bland
WORKDIR /go/src/bland

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download Go modules and dependencies
RUN go mod download



# Manually copy files from Swagger to the correct GOPATH locations
COPY ./Swagger/controller /go/src/bland/controller
COPY ./Swagger/docs /go/src/bland/docs
COPY ./Swagger/docs /go/src/bland/model

# Build the Go application
RUN go build -o /go/bin/app ./main.go

# Set the entry point command to run the Go app
CMD ["/go/bin/app"]

