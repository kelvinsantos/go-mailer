# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
RUN mkdir /app
ADD ./src /app
WORKDIR /app

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o mailer

# Command to run the executable
CMD ["./mailer"]
