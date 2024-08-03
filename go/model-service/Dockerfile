FROM golang:1.22.5-alpine

WORKDIR /app

COPY go/go.mod ./
COPY go/model-service ./model-service/
COPY go/shared ./shared/
COPY sql ./model-service/sql/

RUN go mod download

WORKDIR /app/model-service
RUN go build -o main .

CMD ["./main"]