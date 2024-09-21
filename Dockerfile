# Start with the official Golang image for Go 1.22
FROM golang:1.22-alpine

# Install git and air for live reloading
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@v1.52.3

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose the port the app runs on
EXPOSE 8080

# Build the Go application
# RUN go build -o main .
RUN go build -o main -buildvcs=false .

# Use air to start the app with live reloading
CMD ["air"]
