# Этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем все модули
COPY inventory ./inventory
COPY platform ./platform
COPY shared ./shared

# Собираем приложение без go.work (используем модули напрямую)
WORKDIR /app/inventory
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /inventory-service ./cmd/main.go

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /inventory-service .

# Копируем .env файл для конфигурации приложения
COPY deploy/compose/inventory/.env .env

# Указываем команду запуска
CMD ["./inventory-service"]
