# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
RUN mkdir /app
ADD . /app
RUN cd /app/src && go mod download
WORKDIR /app/src

RUN go build -o gomailer main.go 

EXPOSE 9090

CMD ["./gomailer"]
