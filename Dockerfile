# Этап 1: Сборка приложения
FROM golang:1.22 AS builder

# Установка рабочей директории
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY .env .

# Копируем все файлы приложения
COPY backend .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

# Этап 2: Создание финального образа
FROM alpine:latest

# Копируем собранное приложение из этапа сборки
COPY --from=builder /app/myapp /usr/local/bin/tenderapp
COPY --from=builder /app/.env /usr/local/.env

# Указываем команду для запуска приложения
CMD ["tenderapp"]

# Открываем порт, если необходимо
EXPOSE 8080
