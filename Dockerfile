# syntax=docker/dockerfile:1

FROM golang:1.20
# Use the official Go image as the base image

# Set the working directory in the container
WORKDIR /app

# Copy the application files into the working directory
COPY . /app

RUN go mod download

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 3000

# Define the entry point for the container
CMD ["./main"]