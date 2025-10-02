# Этап, на котором выполняется сборка приложения
FROM golang:1.25-alpine as builder
RUN apk update && apk add --no-cache git && apk upgrade
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/backend-service-profreport app/main.go

# Финальный этап, копируем собранное приложение
FROM alpine:latest

RUN apk update && apk add --no-cache git && apk upgrade

# Аргумент для пути конфигурации
# ARG CONFIG_PATH=/etc/backend-service-profreport
# ARG CONFIG_FILEPATH=config/
# ARG CONFIG_FILENAME=config.yml

# ENV CONFIG_FILEPATH=${CONFIG_FILEPATH}
# ENV CONFIG_FILENAME=${CONFIG_FILENAME}

# Создаем volume для конфигурационных файлов
# VOLUME ${CONFIG_FILEPATH}

COPY --from=builder /build/bin/backend-service-profreport /bin/backend-service-profreport 
# COPY ./config/${CONFIG_FILENAME} ${CONFIG_FILEPATH}${CONFIG_FILENAME}

# RUN mkdir -p /etc/backend-service-profreport/config

ENV HTTP_SERVER_PORT=8080
ENV HTTP_SERVER_ADDRESS_LISTEN=0.0.0.0

EXPOSE $HTTP_SERVER_PORT
ENTRYPOINT ["/bin/backend-service-profreport"]