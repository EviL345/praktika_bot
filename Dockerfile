# Используем многоэтапную сборку для минимизации размера итогового образа
FROM golang:1.24-alpine AS builder

# Установка необходимых зависимостей
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum* ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Убедимся, что директория migrations существует
RUN mkdir -p /app/migrations

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Добавляем поддержку сертификатов и временную зону
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /bot /app/bot
# Копируем конфигурационный файл
COPY config.yaml /app/config.yaml
# Копируем директорию миграций
COPY --from=builder /app/migrations /app/migrations

# Запуск приложения
CMD ["/app/bot"]