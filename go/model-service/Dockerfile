FROM golang:1.23.0-alpine

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

RUN mkdir -p ./config
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -o main .

CMD ["./main"]