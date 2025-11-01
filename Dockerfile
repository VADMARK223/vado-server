# Этап 1: Сборка (builder)
FROM golang:1.25 AS builder

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum - чтобы кэшировался слой зависимостей
COPY go.mod go.sum ./
# Качаем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Сборка бинарника (статическая, без CGO)
# CGO_ENABLED=0 компилято выключает исользование С, и Go собирает чистый статический бинарник. Если чистое CLI, для GUI может все сломать
# -trimpath убирает пути из бинаря (безопасность + меньше размер)
# -o указывает название выходного бинарника
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o vado-server ./cmd/server

# Этап 2: рантайм (минимальный финальный образ)
FROM debian:bookworm-slim AS runtime

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/vado-server .
# Копируем шаблоны и статику
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web/static ./web/static

# Порт для gRPC и HTTP
EXPOSE 50051 5556 8090

# Задаём переменные окружения по умолчанию
ENV PORT=5556 \
    GRPC_PORT=50051 \
    GIN_MODE=release

# Запускаем сервер
CMD ["./vado-server"]

