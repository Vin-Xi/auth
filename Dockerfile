# Stage 1: The builder stage
# Use an official Go image to build the application
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app. CGO_ENABLED=0 creates a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main ./cmd/

# ---

# Stage 2: The final/production stage
# Use a minimal base image for a small and secure final image
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# This command will be run when the container starts
CMD ["./main"]