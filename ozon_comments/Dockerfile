# Базовый образ с минимальной оболочкой alpine
FROM golang:1.22-alpine AS builder

# Установка рабочего каталога
WORKDIR /app

# Установка bash в образ
RUN apk --no-cache add bash

# Загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Сборка приложения
COPY . .
RUN go build -o ./bin/app cmd/main.go

# Копирование рабочих файлов из сборщика
FROM alpine AS worker
COPY --from=builder /app/bin/app /

# Запуск приложения
CMD ["/app"]