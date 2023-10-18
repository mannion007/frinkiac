# Use an official Go runtime as a parent image
FROM golang:1.16

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o frinkiac

# Command to run your application
ENTRYPOINT ["./frinkiac"]