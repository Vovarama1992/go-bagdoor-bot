FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot ./cmd/app

CMD ["./bot"]
