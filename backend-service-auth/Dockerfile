# Этап, на котором выполняется сборка приложения
FROM golang:1.21-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/app-auth-service app/main.go
# Финальный этап, копируем собранное приложение
FROM alpine:latest
COPY --from=builder /build/bin/app-auth-service /bin/app-auth-service
COPY --from=builder /build/config/config_prod.yaml /etc/backend-service-auth/config/config.yaml

ENV HTTP_PORT_EXPOSE=8011
EXPOSE $HTTP_PORT_EXPOSE

ENV CONFIG_PATH=/etc/backend-service-auth/config/config.yaml
ENTRYPOINT ["/bin/app-auth-service"]