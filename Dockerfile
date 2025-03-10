FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . ./

RUN go mod tidy

WORKDIR /app/cmd

RUN go build -o /app/myapp


FROM ubuntu:latest

WORKDIR /root

COPY --from=builder /app .

COPY .env .
COPY mgr.sql .

EXPOSE 8080

CMD ["./myapp"]