FROM golang:1.17.8-alpine3.15 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.io,direct

RUN go build -o main main.go


FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080