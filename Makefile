PROJECT_NAME = vado-app
COMPOSE = docker compose -p $(PROJECT_NAME)
COMPOSE_FULL = $(COMPOSE) -f docker-compose.yml -f docker-compose.kafka.yml

all-up:
	docker compose -p $(PROJECT_NAME) -f docker-compose.yml -f docker-compose.kafka.yml up -d

all-down:
	docker compose -p $(PROJECT_NAME) down

kafka-up:
	$(COMPOSE) -f docker-compose.kafka.yml up -d

kafka-down:
	$(COMPOSE) -f docker-compose.kafka.yml down

up-main:
	docker compose -f docker-compose.yml up -d

down-main:
	docker compose -f docker-compose.yml down

ps:
	$(COMPOSE) ps --format 'table {{.Name}}\t{{.Ports}}'

logs:
	docker compose -p $(PROJECT_NAME) logs vado-server

rebuild:
	$(COMPOSE_FULL) down --volumes
	$(COMPOSE_FULL) up -d --build --remove-orphans

rebuild-full:
	docker compose -p $(PROJECT_NAME) down --rmi all --volumes
	docker compose -p $(PROJECT_NAME) -f docker-compose.yml -f docker-compose.kafka.yml up -d --build

rebuild-server:
	docker compose up -d --build --no-deps vado-server

psql:
	docker exec -it vado-postgres psql -U vadmark -d vadodb

clean:
	docker system prune -af --volumes

PROTO_DIR = api/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
PROTOC = protoc

proto-go:
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

.PHONY: all proto-js clean

all: proto-js

proto-js:
	@echo "Generating gRPC-Web JS files..."
	@mkdir -p $(PB_WEB_OUT_DIR)
	@for file in $(PROTO_FILES); do \
		echo "  -> Compilation $$file"; \
		$(PROTOC) -I=$(PROTO_DIR) $$file \
			--js_out=import_style=commonjs,binary:$(PB_WEB_OUT_DIR) \
			--plugin=protoc-gen-grpc-web=$(GRPC_WEB_PLUGIN) \
			--grpc-web_out=import_style=commonjs,mode=grpcwebtext:$(PB_WEB_OUT_DIR); \
	done
	@echo "Generation complete. Files in $(PB_WEB_OUT_DIR)"

proto-js-clean:
	@echo "Clear $(PB_WEB_OUT_DIR)..."
	rm -rf $(PB_WEB_OUT_DIR)/*.js

bundle:
	npx esbuild web/static/js/grpc.js --bundle --format=esm --outfile=web/static/js/bundle.js


YELLOW:= \033[1;33m
GREEN := \033[1;32m
BLUE  := \033[1;34m
CYAN  := \033[1;36m
RESET := \033[0m

help:
	@echo "$(YELLOW)Available commands:$(RESET)"
	@echo "  $(GREEN)make all-ud$(RESET)          - start all containers"
	@echo "  $(GREEN)make all-down$(RESET)        - stop all containers"
	@echo "  $(GREEN)make kafka-up$(RESET)        - start kafka and kafka UI containers"
	@echo "  $(GREEN)make kafka-down$(RESET)      - stop kafka and kafka UI containers"
	@echo "  $(GREEN)make down-main$(RESET)       - stop server and postgres containers"
	@echo "  $(GREEN)make up-main$(RESET)         - start server and postgres containers"
	@echo "  $(GREEN)make logs$(RESET)            - show logs"
	@echo "  $(GREEN)make rebuild$(RESET)         - rebuild"
	@echo "  $(GREEN)make rebuild-full$(RESET)    - rebuild everything (fresh DB)"
	@echo "  $(GREEN)make rebuild-server$(RESET)  - rebuild only Go server"
	@echo "  $(GREEN)make psql$(RESET)            - open psql shell"
	@echo "  $(GREEN)make clean$(RESET)           - clean Docker cache"
	@echo "  $(GREEN)make proto-go$(RESET)        - generating gRPC Go files"
	@echo "$(CYAN)Web commands:$(RESET)"
	@echo "  $(GREEN)make bundle$(RESET)          - create bundle.js"
	@echo "  $(GREEN)make proto-js$(RESET)        - generating gRPC-Web JS files"
	@echo "  $(GREEN)make proto-js-clean$(RESET)  - remove gRPC-Web JS files"
.DEFAULT_GOAL := help