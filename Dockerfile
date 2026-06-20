# ---- Build Stage ----
FROM golang:1.25.8-alpine3.23 AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Install go dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy only necessary source files
COPY *.go ./
COPY config/ ./config/
COPY controllers/ ./controllers/
COPY middleware/ ./middleware/
COPY routes/ ./routes/
COPY utils/ ./utils/

# Build the Go application
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -a -installsuffix cgo -ldflags="-w -s" -o main main.go

# ---- Final Stage ----
FROM alpine:3.21
WORKDIR /app

# Install minimal runtime dependencies
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the images directory (only if needed at runtime)
COPY images/ ./images/

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8090

# Run the compiled binary
CMD ["./main"]