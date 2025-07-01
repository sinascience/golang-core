# ---- Build Stage ----
# This stage has the Go toolchain and is used for development.
FROM golang:1.24-alpine AS builder

# Set the working directory to /app
WORKDIR /app

# Copy and download dependencies first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Install Air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy the rest of the source code
COPY . .


# ---- Final Stage ----
# This stage is for creating a minimal production image without the Go toolchain.
FROM alpine:latest

# Copy only the compiled application from the builder stage
# (This assumes a production build step would place the binary here)
COPY --from=builder /app/tmp/server /usr/local/bin/server

# Expose the port
EXPOSE 3000

# Set the default command for production
CMD ["server"]