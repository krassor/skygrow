# Этап, на котором выполняется сборка приложения
FROM golang:1.23-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/app-backend-miniaps-bot app/main.go
# Финальный этап, копируем собранное приложение
FROM alpine:latest
COPY --from=builder /build/bin/app-backend-miniaps-bot  /bin/app-backend-miniaps-bot
# ENV DEVICES_DB_NAME=postgres
# ENV DEVICES_DB_PASSWORD=postgres
# ENV DEVICES_DB_USER=postgres
# ENV DEVICES_DB_HOST=172.17.0.2
# ENV DEVICES_DB_PORT=5432
ENV HTTP_SERVER_PORT=8080
ENV HTTP_SERVER_ADDRESS_LISTEN=0.0.0.0
ENV TGBOT_APITOKEN=5995322288:AAE9eClgD8cQyasmv5aywzVSaZ5yu-zcjcU
ENV USER_MESSAGE_DB_TYPE=Inmemory
ENV CONFIG_FILEPATH=/etc/backend-miniaps-bot/config.yml
# EXPOSE $HTTP_SERVER_PORT
ENTRYPOINT ["/bin/app-backend-miniaps-bot"]