FROM golang:1.23-alpine3.20 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/bot cmd/bot/main.go

FROM alpine:3.20

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/bot /app/bot
COPY --from=builder /app/config/config-bot.yaml /app/config-bot.yaml

ENTRYPOINT ["./bot", "-config", "config-bot.yaml"]