FROM golang:1.23-alpine AS builder

COPY . /github.com/codeboris/pocket_telegram_bot/
WORKDIR /github.com/codeboris/pocket_telegram_bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bo/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/codeboris/pocket_telegram_bot/bin/bot .
COPY --from=builder /github.com/codeboris/pocket_telegram_bot/configs configs/

EXPOSE 80

CMD [ "./bot" ]