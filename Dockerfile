# Use the official Golang image from the Docker Hub
FROM golang:1.19-alpine

# Install air for live reloading
RUN apk add --no-cache git
RUN go install github.com/cosmtrek/air@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Expose port 7099 to the outside world
EXPOSE 7099

# Run the application with air for live reloading
CMD ["air"]
