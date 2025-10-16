# gRPC

Генерация из `.proto` файла

```shell
protoc --go_out=./ --go-grpc_out=./ api/proto/hello.proto
```

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

# Golang

Инициализация проекта
```shell
go mod init vado_server
```
Установка `zap`
```shell
go get -u go.uber.org/zap
```
Запуск проекта
```shell
go run ./cmd/server/main.go
```

# Stack
- gRPC
- Postgres
- Zap
- Gin
- Gorm
- JWT

На будущее: golang-migrate