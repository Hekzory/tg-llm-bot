FROM golang:1.24.1-alpine

# Add non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy go.mod and go.sum files first
COPY go/go.mod ./
COPY go/telegram-service/go.mod go/telegram-service/go.sum ./telegram-service/
COPY go/shared/go.mod go/shared/go.sum ./shared/

# Download dependencies
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

COPY go/telegram-service ./telegram-service/
COPY go/shared ./shared/
COPY sql ./telegram-service/sql/

WORKDIR /app/telegram-service

# Create config directory with proper permissions
RUN mkdir -p /app/telegram-service/config && \
    chown -R appuser:appgroup /app/telegram-service/config

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -o main .

# Change ownership of the application directory
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

CMD ["./main"]
