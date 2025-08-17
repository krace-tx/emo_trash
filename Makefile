# 工具检查
GO := go
GOZERO := goctl
MKDIR := mkdir -p
RM := rm -rf
TOUCH := touch
ECHO := echo -e
SED := sed

# 项目配置
PROJECT_NAME := emo_trash
PROJECT_PATH := github.com/krace-tx/emo_trash

# 颜色定义
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

# 目录配置
API_SRC_DIR := cmd/api
RPC_SRC_DIR := cmd/rpc
API_DEST_DIR := app/api
RPC_DEST_DIR := app/rpc
SWAGGER_DEST_DIR := swagger
DOCKER_COMPOSE_DIR := deploy

# 查找文件
API_FILES := $(wildcard $(API_SRC_DIR)/*.api)
PROTO_FILES := $(wildcard $(RPC_SRC_DIR)/*.proto)

# 生成目标
API_TARGETS := $(patsubst $(API_SRC_DIR)/%.api,$(API_DEST_DIR)/%/.,$(API_FILES))
RPC_TARGETS := $(patsubst $(RPC_SRC_DIR)/%.proto,$(RPC_DEST_DIR)/%/.,$(PROTO_FILES))

# 默认目标
.DEFAULT_GOAL := help

.PHONY: all init api rpc check mod clean swagger docker-compose wsl help

all: api rpc
	@$(ECHO) "$(GREEN)All code generated successfully!$(NC)"

# 初始化项目环境
init:
	@$(ECHO) "$(YELLOW)Initializing project environment...$(NC)"

	# 创建必要的源目录
	@$(MKDIR) $(API_SRC_DIR) $(RPC_SRC_DIR) $(API_DEST_DIR) $(RPC_DEST_DIR)
	@$(ECHO) "$(GREEN)Created required directories:$(NC) $(API_SRC_DIR), $(RPC_SRC_DIR), $(API_DEST_DIR), $(RPC_DEST_DIR)"

	# 安装Go依赖工具
	@$(ECHO) "$(YELLOW)Installing required Go tools...$(NC)"
	@go install github.com/zeromicro/go-zero/tools/goctl@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest
	# 检查GOPATH配置
	@echo 'export PATH=$PATH:/go/bin' >> ~/.bashrc
	@source ~/.bashrc

	# 检查goctl是否可以执行
	@goctl env check --install --verbose --force || (echo "$(RED)Error: goctl is not found in PATH$(NC) && exit 1")

# API代码生成
api: $(API_TARGETS)

$(API_DEST_DIR)/%/. : $(API_SRC_DIR)/%.api
	@$(ECHO) "$(YELLOW)Generating API code for $<...$(NC)"
	@# 捕获命令输出和返回码，支持跳过 "missing service" 错误
	@OUTPUT=$$($(GOZERO) api go -api $< -dir $(API_DEST_DIR)/$(notdir $(basename $<)) 2>&1); \
	RET_CODE=$$?; \
	if [ $$RET_CODE -eq 0 ]; then \
		$(TOUCH) $@; \
		$(ECHO) "$(GREEN)API code generated to $(API_DEST_DIR)/$(notdir $(basename $<))$(NC)"; \
	elif echo "$$OUTPUT" | grep -q "Error: missing service"; then \
		$(ECHO) "$(YELLOW)Skipping $<: missing service definition$(NC)"; \
	else \
		$(ECHO) "$(RED)Failed to generate API code for $<:$$OUTPUT$(NC)"; \
		exit 1; \
	fi

# RPC代码生成
rpc: $(RPC_TARGETS)

$(RPC_DEST_DIR)/%/. : $(RPC_SRC_DIR)/%.proto
	@PROTOC_OUTPUT=$$(cd $(RPC_SRC_DIR) && \
	protoc $(notdir $<) \
	--proto_path=. \
	--proto_path=../../third_party \
	--validate_out=paths=source_relative,lang=go,template=validator_zh.yaml:../../$(RPC_DEST_DIR)/$(notdir $(basename $<))/pb \
	 2>&1); \
	PROTOC_RET=$$?; \
	if [ $$PROTOC_RET -ne 0 ]; then \
		$(ECHO) "$(RED)Failed to generate protobuf code: $$PROTOC_OUTPUT$(NC)"; \
		exit 1; \
	fi
	@$(ECHO) "$(YELLOW)Generating RPC code for $<...$(NC)"
	@# 捕获命令输出和返回码，支持跳过 "rpc service not found" 错误
	@OUTPUT=$$(cd $(RPC_SRC_DIR) && \
	 $(GOZERO) rpc protoc $(notdir $<) \
  	 --proto_path=. \
  	 --proto_path=../../third_party \
	 --go_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --go-grpc_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --zrpc_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --style=go_zero \
	 -m 2>&1); \
	RET_CODE=$$?; \
	if [ $$RET_CODE -eq 0 ]; then \
		$(ECHO) "$(GREEN)RPC code generated to $(RPC_DEST_DIR)/$(notdir $(basename $<))$(NC)"; \
		$(ECHO) "Injecting custom tags to remove omitempty..."; \
		find $(RPC_DEST_DIR)/$(notdir $(basename $<))/pb -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;; \
		$(TOUCH) $@; \
	elif echo "$$OUTPUT" | grep -q "rpc service not found"; then \
		$(ECHO) "$(YELLOW)Skipping $<: rpc service not found$(NC)"; \
	else \
		$(ECHO) "$(RED)Failed to generate RPC code for $<:$$OUTPUT$(NC)"; \
		exit 1; \
	fi

# 检查文件
check:
	@$(ECHO) "$(YELLOW)API files found:$(NC) $(API_FILES)"
	@$(ECHO) "$(YELLOW)Proto files found:$(NC) $(PROTO_FILES)"

mod:
	rm go.mod go.sum
	go mod init $(PROJECT_PATH) && go mod tidy

clean:
	@git update-index --assume-unchanged go.mod
	@git update-index --assume-unchanged go.sum

swagger:
	@$(ECHO) "$(YELLOW)Generating Swagger documentation...$(NC)"
	@$(GOZERO) api swagger -api $(API_FILES) -dir $(SWAGGER_DEST_DIR)

docker-compose:
	@$(ECHO) "$(YELLOW)=== Starting Docker Compose Services ==="$(NC)

	# Check if docker-compose (V1 or V2) is installed
	@if ! command -v docker-compose &> /dev/null && ! command -v docker compose &> /dev/null; then \
		$(ECHO) "$(RED)Error: docker-compose is not installed. Please install it first.$(NC)"; \
		exit 1; \
	fi

	# Check if docker-compose.yaml exists
	@if [ ! -f "$(DOCKER_COMPOSE_DIR)/docker-compose.yaml" ]; then \
		$(ECHO) "$(RED)Error: docker-compose.yaml not found in $(DOCKER_COMPOSE_DIR)/$(NC)"; \
		exit 1; \
	fi

	@cd $(DOCKER_COMPOSE_DIR) && \
	  $(ECHO) "$(YELLOW)=== Stopping existing containers ==="$(NC) && \
	  # Use docker-compose if available, otherwise docker compose
	  if command -v docker-compose &> /dev/null; then \
		docker-compose down && \
		docker-compose pull && \
		docker-compose up -d --build --force-recreate; \
	  else \
		docker compose down && \
		docker compose pull && \
		docker compose up -d --build --force-recreate; \
	  fi && \
	  $(ECHO) "$(GREEN)=== Docker Compose services started successfully ==="$(NC)

wsl:
	docker exec -it emo_trash_dev bash

# 帮助信息
help:
	@$(ECHO) "$(GREEN)Available targets:$(NC)"
	@$(ECHO) "  $(YELLOW)all$(NC)        - Generate all API and RPC code (default)"
	@$(ECHO) "  $(YELLOW)api$(NC)        - Generate API code"
	@$(ECHO) "  $(YELLOW)rpc$(NC)        - Generate RPC code"
	@$(ECHO) "  $(YELLOW)check$(NC)      - Check for API and Proto files"
	@$(ECHO) "  $(YELLOW)init$(NC)       - Initialize project environment (install tools, create dirs)"
	@$(ECHO) "  $(YELLOW)mod$(NC)        - Initialize project environment (mod init)"
	@$(ECHO) "  $(YELLOW)clean$(NC)      - Remove git gitignore files"
	@$(ECHO) "  $(YELLOW)help$(NC)       - Show this help message"
	@$(ECHO) ""
	@$(ECHO) "$(GREEN)Examples:$(NC)"
	@$(ECHO) "  make init           # Initialize project environment"
	@$(ECHO) "  make                # Generate all code"
	@$(ECHO) "  make rpc            # Generate only RPC code"
	@$(ECHO) "  make check          # Check for API and Proto files"
	@$(ECHO) "  make swagger        # Generate Swagger documentation"
	@$(ECHO) "  make docker-compose # Start Docker Compose services"
	@$(ECHO) "  make help           # Show this help message"


