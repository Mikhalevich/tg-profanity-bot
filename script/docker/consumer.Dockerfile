FROM golang:1.22-alpine3.19 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/consumer cmd/consumer/main.go

FROM alpine:3.19

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/consumer /app/consumer
COPY --from=builder /app/config/config-consumer.yaml /app/config-consumer.yaml

ENTRYPOINT ["./consumer", "-config", "config-consumer.yaml"]