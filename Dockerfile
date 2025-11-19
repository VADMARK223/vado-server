# Этап 1: Сборка (builder) слой с зависимостями
FROM golang:1.25 AS builder

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum - чтобы кэшировался слой зависимостей
COPY go.mod go.sum ./
# Качаем зависимости
RUN go mod download

# Копируем только исходники
COPY api ./api
COPY "cmd" "./cmd"
COPY internal ./internal
COPY web ./web

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
COPY --from=builder /app/web ./web

# Порт для gRPC и HTTP
EXPOSE 50051 5555 8090

# Устанавливаем базовые env (не мешают docker-compose)
ENV APP_ENV=production \
    GIN_MODE=release

# Запускаем сервер
CMD ["./vado-server"]