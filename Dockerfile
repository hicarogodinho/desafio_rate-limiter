FROM golang:1.20-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o rate-limite-go .cmd/server

EXPOSE 8080

CMD ["./desfio_rate-limiter"]