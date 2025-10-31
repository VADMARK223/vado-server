PROJECT_NAME = vado-app

all-up:
	docker compose -p $(PROJECT_NAME) -f docker-compose.yml -f docker-compose.kafka.yml up -d

all-down:
	docker compose -p $(PROJECT_NAME) down

kafka-up:
	docker compose -f docker-compose.kafka.yml up -d

kafka-down:
	docker compose -f docker-compose.kafka.yml down

up-main:
	docker compose -f docker-compose.yml up -d

down-main:
	docker compose -f docker-compose.yml down

ps:
	docker compose -p $(PROJECT_NAME) ps vado-server

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

PROTO_DIR = api/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
PROTOC = protoc

go-proto:
	@echo "Generating Go gRPC files..."
	@for file in $(PROTO_FILES); do \
		echo "  -> Compilation $$file"; \
		$(PROTOC) -I=$(PROTO_DIR) $$file \
			--go_out=. \
			--go-grpc_out=. ; \
	done
	@echo "Generation complete."

PB_WEB_OUT_DIR = web/static/js/pb
GRPC_WEB_PLUGIN = $(shell which protoc-gen-grpc-web)

.PHONY: all web-proto clean

all: web-proto

web-proto:
	@echo "Generating gRPC-Web JS files..."
	@mkdir -p $(PB_WEB_OUT_DIR)
	@for file in $(PROTO_FILES); do \
		echo "  -> Compilation $$file"; \
		$(PROTOC) -I=$(PROTO_DIR) $$file \
			--js_out=import_style=closure,binary:$(PB_WEB_OUT_DIR) \
			--plugin=protoc-gen-grpc-web=$(GRPC_WEB_PLUGIN) \
			--grpc-web_out=import_style=closure,mode=grpcweb:$(PB_WEB_OUT_DIR); \
	done
	@echo "Generation complete. Files in $(PB_WEB_OUT_DIR)"

web-proto-clean:
	@echo "Clear $(PB_WEB_OUT_DIR)..."
	rm -rf $(PB_WEB_OUT_DIR)/*.js


YELLOW := \033[1;33m
GREEN := \033[1;32m
RESET := \033[0m

help:
	@echo "$(YELLOW)Available command:$(RESET)"
	@echo "  $(GREEN)make all-ud$(RESET)          - start all containers"
	@echo "  $(GREEN)make all-down$(RESET)        - stop all containers"
	@echo "  $(GREEN)make kafka-up$(RESET)        - start kafka and kafka UI containers"
	@echo "  $(GREEN)make kafka-down$(RESET)      - stop kafka and kafka UI containers"
	@echo "  $(GREEN)make down-main$(RESET)       - stop server and postgres containers"
	@echo "  $(GREEN)make up-main$(RESET)         - start server and postgres containers"
	@echo "  $(GREEN)make logs$(RESET)            - show logs"
	@echo "  $(GREEN)make rebuild$(RESET)         - rebuild everything (fresh DB)"
	@echo "  $(GREEN)make rebuild-server$(RESET)  - rebuild only Go server"
	@echo "  $(GREEN)make psql$(RESET)            - open psql shell"
	@echo "  $(GREEN)make clean$(RESET)           - clean Docker cache"
	@echo "  $(GREEN)make web-proto$(RESET)       - generating gRPC-Web JS files"
	@echo "  $(GREEN)make web-proto-clean$(RESET) - remove gRPC-Web JS files"
	@echo "  $(GREEN)make go-proto$(RESET)        - generating gRPC files"
.DEFAULT_GOAL := help