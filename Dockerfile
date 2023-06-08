FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN go build -o server main.go

CMD ["./server"]
