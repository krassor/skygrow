# Этап, на котором выполняется сборка приложения
FROM golang:1.20-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/app-auth-service app/main.go
# Финальный этап, копируем собранное приложение
FROM alpine:latest
COPY --from=builder /build/bin/app-auth-service /bin/app-auth-service
# ENV DB_NAME=skygrowdb
# ENV DB_PASSWORD=copro123@
# ENV DB_USER=krassor
# ENV DB_HOST=skygrow-pg
# ENV DB_PORT=5432
# ENV HTTP_PORT=8011
# ENV HTTP_HOST_LISTEN=0.0.0.0
# ENV TGBOT_APITOKEN=5817824261:AAEoDtCT--VWk8ESKNTBftndBnQ6MgdtJxg
# ENV HTTP_PORT_EXPOSE=8011
# EXPOSE $HTTP_PORT_EXPOSE
ENTRYPOINT ["/bin/app-auth-service"]