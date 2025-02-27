# Этап, на котором выполняется сборка приложения
FROM golang:1.23-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/backend-miniaps-bot app/main.go

# Финальный этап, копируем собранное приложение
FROM alpine:latest

# Аргумент для пути конфигурации
# ARG CONFIG_PATH=/etc/backend-miniaps-bot
ARG CONFIG_FILENAME=config.yml
ARG CONFIG_FILEPATH=config/
ENV CONFIG_FILEPATH=${CONFIG_FILEPATH}
ENV CONFIG_FILENAME=${CONFIG_FILENAME}

COPY --from=builder /build/bin/backend-miniaps-bot /bin/backend-miniaps-bot 
COPY ./config/${CONFIG_FILENAME} ${CONFIG_FILEPATH}${CONFIG_FILENAME}

ENV HTTP_SERVER_PORT=8080
ENV HTTP_SERVER_ADDRESS_LISTEN=0.0.0.0


# Создаем volume для конфигурационных файлов
VOLUME ${CONFIG_FILEPATH}

EXPOSE $HTTP_SERVER_PORT
ENTRYPOINT ["/bin/backend-miniaps-bot"]