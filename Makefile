PROJECT_NAME = vado-app

up-all:
	docker compose -p $(PROJECT_NAME) -f docker-compose.yml -f docker-compose.kafka.yml up -d

down-all:
	docker compose -p $(PROJECT_NAME) down

up-main:
	docker compose -f docker-compose.yml up -d

down-main:
	docker compose -f docker-compose.yml down

up-kafka:
	docker compose -f docker-compose.kafka.yml up -d

down-kafka:
	docker compose -f docker-compose.kafka.yml down

ps:
	docker compose -p $(PROJECT_NAME) ps -f vado-server

logs:
	docker compose -p $(PROJECT_NAME) logs vado-server

# Полный rebuild (удаляет все контейнеры и тома, пересоздаёт базу и запускает init-скрипты)
rebuild:
	docker compose down -v
	docker compose up -d --build

# Пересобрать сервер, не трогая базу
rebuild-server:
	docker compose up -d --build --no-deps vado-server

psql:
	docker exec -it vado-postgres psql -U vadmark -d vadodb

clean:
	docker system prune -af --volumes

YELLOW := \033[1;33m
GREEN := \033[1;32m
RESET := \033[0m

help:
	@echo "$(YELLOW)Available command:$(RESET)"
	@echo "  $(GREEN)make up-all$(RESET)          - start all containers"
	@echo "  $(GREEN)make down-all$(RESET)        - stop all containers"
	@echo "  $(GREEN)make up-main$(RESET)         - start server and postgres containers"
	@echo "  $(GREEN)make down-main$(RESET)       - stop server and postgres containers"
	@echo "  $(GREEN)make up-kafka$(RESET)        - start kafka and kafka UI containers"
	@echo "  $(GREEN)make down-kafka$(RESET)      - stop kafka and kafka UI containers"
	@echo "  $(GREEN)make logs$(RESET)            - show logs"
	@echo "  $(GREEN)make rebuild$(RESET)         - rebuild everything (fresh DB)"
	@echo "  $(GREEN)make rebuild-server$(RESET)  - rebuild only Go server"
	@echo "  $(GREEN)make psql$(RESET)            - open psql shell"
	@echo "  $(GREEN)make clean$(RESET)           - clean Docker cache"
.DEFAULT_GOAL := help