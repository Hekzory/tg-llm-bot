FROM golang:1.24.2-alpine

# Add non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy go.mod and go.sum files first
COPY go/go.mod ./
COPY go/model-service/go.mod go/model-service/go.sum ./model-service/
COPY go/shared/go.mod go/shared/go.sum ./shared/

# Download dependencies
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

COPY go/model-service ./model-service/
COPY go/shared ./shared/
COPY sql ./model-service/sql/

WORKDIR /app/model-service

# Create config directory with proper permissions
RUN mkdir -p /app/model-service/config && \
    chown -R appuser:appgroup /app/model-service/config

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -o main .

# Change ownership of the application directory
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

CMD ["./main"]
