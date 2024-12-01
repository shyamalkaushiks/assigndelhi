
FROM golang:1.22.2 AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

FROM alpine:latest

WORKDIR /root/


COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
