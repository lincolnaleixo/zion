# Use the official Golang image as a base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./main.go

# Expose the port the app runs on
EXPOSE 8080

# Run the executable
CMD ["./main"]
