# =========================
# 🎨 HELP SECTION
# =========================
MAKEFLAGS += --no-print-directory
YELLOW:= \033[1;33m
GREEN := \033[1;32m
BLUE  := \033[1;34m
CYAN  := \033[1;36m
ORANGE := \033[38;5;208m
RESET := \033[0m

# =========================
# Read .env.prod
# =========================
ifneq (,$(wildcard .env.prod))
    include .env.prod
    export $(shell sed -n 's/^\([^#[:space:]]\+\)=.*/\1/p' .env.prod)
endif
ifeq ($(KAFKA_ENABLE), true)
    KAFKA_FILE = -f docker-compose.kafka.yml
else
    KAFKA_FILE =
endif

PROJECT_NAME = vado-app
COMPOSE = docker compose -p $(PROJECT_NAME)
COMPOSE_FULL = $(COMPOSE) -f docker-compose.yml $(KAFKA_YML)

PROTO_DIR = api/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
PROTOC = protoc

all-up:
	docker compose -p $(PROJECT_NAME) -f docker-compose.yml $(KAFKA_YML) up -d

all-down:
	docker compose -p $(PROJECT_NAME) down

kafka-up:
	$(COMPOSE) $(KAFKA_YML) up -d

kafka-down:
	$(COMPOSE) $(KAFKA_YML) down

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
	docker compose -p $(PROJECT_NAME) -f docker-compose.yml $(KAFKA_YML) up -d --build

rebuild-server:
	docker compose up -d --build --no-deps vado-server

psql:
	docker exec -it vado-postgres psql -U vadmark -d vadodb

clean:
	docker system prune -af --volumes

proto-go:
	@echo "Generating Go gRPC files..."
	@for file in $(PROTO_FILES); do \
		echo "  -> Compilation $$file"; \
		$(PROTOC) -I=$(PROTO_DIR) $$file \
			--go_out=. \
			--go-grpc_out=. ; \
	done
	@echo "✅ Generation complete."

PB_WEB_OUT_DIR = ./web/static/js/pb
GRPC_WEB_PLUGIN = $(shell which protoc-gen-grpc-web)
TS_PLUGIN := $(shell pwd)/node_modules/.bin/protoc-gen-ts

proto-ts-clean:
	@echo "$(ORANGE)⚠️ Clear all *.ts$(PB_WEB_OUT_DIR)...$(RESET)"
	@find $(PB_WEB_OUT_DIR) -type f \( -name "*.ts" -o -name "*.js" \) -delete
	@echo "$(GREEN)✅️ Cleaning is complete$(RESET)"

proto-ts:
	@echo "🔧 Generating gRPC-Web TypeScript files..."
	@mkdir -p $(PB_WEB_OUT_DIR)
	@for file in $(PROTO_DIR)/*.proto; do \
		echo "  🔵 Compiling $$file"; \
		protoc -I=$(PROTO_DIR) \
			--plugin=protoc-gen-ts=$(TS_PLUGIN) \
			--js_out=import_style=commonjs,binary:$(PB_WEB_OUT_DIR) \
			--ts_out=service=grpc-web:$(PB_WEB_OUT_DIR) \
			$$file; \
	done
	@echo "✅ TypeScript gRPC stubs generated → $(PB_WEB_OUT_DIR)"

bundle:
	@echo "$(BLUE)📦 Bundling TypeScript client with esbuild...$(RESET)"
	npx esbuild web/static/js/grpc.ts --bundle --format=esm --outfile=web/static/js/bundle.js --platform=browser --target=es2020 --define:process.env.GRPC_WEB_PORT="'$(GRPC_WEB_PORT)'"
	@echo "$(GREEN)✅ Bundle created → web/static/js/bundle.js$(RESET)"

proto-ts-all: ## 🚀 Full pipeline: clean → generate → bundle
	@echo "$(BLUE)🚀 Starting full gRPC-Web TypeScript build pipeline...$(RESET)"
	@$(MAKE) proto-ts-clean || { echo "$(ORANGE)❌ Stage failed: proto-ts-clean$(RESET)"; exit 1; }
	@$(MAKE) proto-ts || { echo "$(ORANGE)❌ Stage failed: proto-ts$(RESET)"; exit 1; }
	@$(MAKE) bundle || { echo "$(ORANGE)❌ Stage failed: bundle$(RESET)"; exit 1; }
	@echo "$(GREEN)✅ All stages completed successfully!$(RESET)"

help:
	@echo "$(YELLOW)🧩 Available Make targets:$(RESET)"
	@echo ""
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
	@echo ""
	@echo "$(CYAN)Type script proto:$(RESET)"
	@echo "  $(GREEN)make proto-ts-clean$(RESET)  - 🧹 Clean generated *.ts and *.js, files from $(PB_WEB_OUT_DIR)"
	@echo "  $(GREEN)make proto-ts$(RESET)        - 🔧 Generate gRPC-Web client files (.js, .d.ts, .ts)"
	@echo "  $(GREEN)make bundle$(RESET)          - 📦 Bundle TypeScript client into a single bundle.js using esbuild"
	@echo "  $(GREEN)make proto-ts-all$(RESET)    - 🚀 Run the full pipeline: clean → generate → bundle"
.DEFAULT_GOAL := help