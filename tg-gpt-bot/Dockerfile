# Этап, на котором выполняется сборка приложения
FROM golang:1.20-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/app-tg-gpt-bot app/main.go
# Финальный этап, копируем собранное приложение
FROM alpine:latest
COPY --from=builder /build/bin/app-tg-gpt-bot  /bin/app-tg-gpt-bot 
# ENV DEVICES_DB_NAME=postgres
# ENV DEVICES_DB_PASSWORD=postgres
# ENV DEVICES_DB_USER=postgres
# ENV DEVICES_DB_HOST=172.17.0.2
# ENV DEVICES_DB_PORT=5432
# ENV DEVICES_HTTP_PORT=8011
# ENV DEVICES_HTTP_HOST_LISTEN=0.0.0.0
ENV TGBOT_APITOKEN=6037249375:AAGXTPvT2zr2LAyty9OtUYEov-oCH15sZ8w
ENV OPENAI_TOKEN=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
ENV USER_MESSAGE_DB_TYPE=Inmemory
ENV CONFIG_FILEPATH=/etc/tg-gpt-bot/config.yml
# EXPOSE $DEVICES_HTTP_PORT
ENTRYPOINT ["/bin/app-tg-gpt-bot"]