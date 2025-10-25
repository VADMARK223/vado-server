# Запустить контейнеры (без пересборки)
up-dc:
	docker compose up -d

# Остановить контейнеры (не удаляя данные)
down-dc:
	docker compose down

# Полный rebuild (удаляет все контейнеры и тома, пересоздаёт базу и запускает init-скрипты)
rebuild:
	docker compose down -v
	docker compose up -d --build

# Пересобрать сервер, не трогая базу
rebuild-server:
	docker compose up -d --build --no-deps vado-server

logs:
	docker compose logs -f vado-server

psql:
	docker exec -it vado-postgres psql -U vadmark -d vadodb

clean:
	docker system prune -af --volumes

YELLOW := \033[1;33m
GREEN := \033[1;32m
RESET := \033[0m

help:
	@echo "$(YELLOW)Available command:$(RESET)"
	@echo "  $(GREEN)make up-dc$(RESET)           - start containers"
	@echo "  $(GREEN)make down-dc$(RESET)         - stop containers"
	@echo "  $(GREEN)make rebuild$(RESET)         - rebuild everything (fresh DB)"
	@echo "  $(GREEN)make rebuild-server$(RESET)  - rebuild only Go server"
	@echo "  $(GREEN)make logs$(RESET)            - show vado-server logs"
	@echo "  $(GREEN)make psql$(RESET)            - open psql shell"
	@echo "  $(GREEN)make clean$(RESET)           - clean Docker cache"
.DEFAULT_GOAL := help