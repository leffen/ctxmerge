# Stage 1: Build the application
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/ctxmerge main.go

# Stage 2: Create a minimal distroless image
FROM gcr.io/distroless/static:nonroot

# Set the working directory in the distroless image
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/ctxmerge /app/ctxmerge

# Use a non-root user
USER nonroot:nonroot

# Command to run the application
ENTRYPOINT ["/app/ctxmerge"]
