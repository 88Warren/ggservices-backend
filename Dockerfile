# Multi-arch compatible — builds for linux/amd64 and linux/arm64
# docker buildx build \
#   --platform linux/amd64,linux/arm64 \
#   -t lmwcode/ggservices-backend:v1.0.0 \
#   --no-cache --push .

# ── Stage 1: build ────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

# Install go dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY *.go ./
COPY config/ ./config/
COPY controllers/ ./controllers/
COPY middleware/ ./middleware/
COPY routes/ ./routes/
COPY utils/ ./utils/

# CGO_ENABLED=0 = pure Go binary, no C compiler needed
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -a -ldflags="-w -s" -o main main.go

# ── Stage 2: run ──────────────────────────────────────────────
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/main .

# Copy images served as static files
COPY images/ ./images/

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8090

CMD ["./main"]
