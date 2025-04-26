# Первый Dockerfile
FROM golang:1.23.6 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальной код проекта
COPY . .

# Устанавливаем нужную версию swag
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

# Генерируем swagger-документацию
RUN swag init

# Собираем бинарный файл
RUN go build -o main main.go

# Используем минимальный образ для запуска
FROM ubuntu:latest

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию.
WORKDIR /root/

# Создаем директорию для хранения конфигов
RUN mkdir "configs"

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/main .

# Копируем конфиги
COPY --from=builder /app/configs/docker/configs.json ./configs
COPY --from=builder /app/configs/docker/example.json ./configs

# Копируем переменные окружения
COPY --from=builder /app/.env .
COPY --from=builder /app/example.env .

# Если Swagger-файлы нужны на runtime, их тоже можно скопировать (опционально):
# COPY --from=builder /app/docs ./docs

# Открываем порт
EXPOSE 6565

# Команда для запуска
CMD ["./main"]
