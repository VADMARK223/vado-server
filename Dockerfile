# Этап сборки
FROM golang:1.25 AS builder

WORKDIR /app

# Копируем go.mod и go.sum (Кэшируем зависиомсти)
COPY go.mod go.sum ./

# Качаем зависимости
RUN go mod download

# Копируем исходики
COPY . .

# Собираем бинарник
# * CGO_ENABLED=0 компилято выключает исользование С, и Go собирает чистый статический бинарник. Если чистое CLI, для GUI может все сломать
# * -o указывает название выходного бинарника
RUN CGO_ENABLED=0 go build -o vado-ping ./cmd/server

# Запуск этапа (минимальный финальный образ)
FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/vado-server .
# Копируем шаблоны и статику
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web/static ./web/static

# Пробрасываем порт
EXPOSE 5556

# Задаём переменные окружения по умолчанию
ENV PORT=5556 \
    GRPC_PORT=50051 \
    GIN_MODE=release

# Запускаем сервер
CMD ["./vado-server"]

