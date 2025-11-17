# gRPC

Генерация из `.proto` файла

```shell
protoc --go_out=./ --go-grpc_out=./ api/proto/hello.proto
```
## gRPC-curl
```shell
grpcurl -plaintext -import-path ./api/proto -proto auth.proto -d '{"username": "1", "password": "1"}' localhost:50051 AuthService/Login
grpcurl -plaintext -import-path ./api/proto -proto ping.proto localhost:50051 PingService/Ping
````

# Linux

Прибить порт
```shell
sudo lsof -i:8080
sudo kill -9 PID
```
# Docker
## Postgres

Зайти в контейнер
```shell
docker exec -it vado_postgres bash
```

### psql
Зайти в `psql`
```shell
psql -U vadmark -d vadodb
```
Список таблиц
```shell
\dt
```
Структура таблицы
```shell
\d tasks
```
Удаление volume, чтобы база создалась заново
```shell
docker volume rm vado-server_postgres-data
```

# Golang

Инициализация проекта
```shell
go mod init vado_server
```
Чистит зависимости (добавляет/удаляет лишние)
```shell
go mod tidy
```
Обновить все зависимости:
```shell
go get -u ./...
go mod tidy

```

Показывает, почему модуль был добавлен
```shell
go mod why <package>
```

Установка `zap`
```shell
go get -u go.uber.org/zap
```

Показывает все модули
```shell
go list -m all
```

Запуск проекта
```shell
go run ./cmd/ping/main.go
```

# Stack
- gRPC
- gRPC Web
- Postgres
- Zap
- Gin
- Gorm
- JWT (access + refresh)

На будущее: golang-migrate