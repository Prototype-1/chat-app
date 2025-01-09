# Use the official Go image as the base image for building the app
FROM golang:1.22.1 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Build the application with static linking (CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory in the production container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]


