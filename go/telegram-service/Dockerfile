FROM golang:1.22.5-alpine

WORKDIR /app

COPY go/go.mod ./
COPY go/telegram-service ./telegram-service/
COPY go/shared ./shared/
COPY sql ./telegram-service/sql/

RUN go mod download

WORKDIR /app/telegram-service
RUN go build -o main .

CMD ["./main"]