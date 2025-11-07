# lightweight runtime-only image
FROM debian:bookworm-slim

WORKDIR /app

# копируем бинарь внутрь контейнера
COPY vado-server .

# открываем нужные порты
EXPOSE 8090 50051 5556

# запускаем
CMD ["./vado-server"]