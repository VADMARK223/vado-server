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
- Postgres
- Zap
- Gin
- Gorm

На будущее: golang-migrate