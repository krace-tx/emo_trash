# 工具检查
GO := go
GOZERO := goctl
MKDIR := mkdir -p
RM := rm -rf
TOUCH := touch
ECHO := echo
SED := sed

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

# 查找文件
API_FILES := $(wildcard $(API_SRC_DIR)/*.api)
PROTO_FILES := $(wildcard $(RPC_SRC_DIR)/*.proto)

# 生成目标
API_TARGETS := $(patsubst $(API_SRC_DIR)/%.api,$(API_DEST_DIR)/%/.,$(API_FILES))
RPC_TARGETS := $(patsubst $(RPC_SRC_DIR)/%.proto,$(RPC_DEST_DIR)/%/.,$(PROTO_FILES))

# 默认目标
.DEFAULT_GOAL := help

.PHONY: all init api rpc check help

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
	@$(ECHO) "$(YELLOW)Verifying GOPATH configuration...$(NC)"
	@if [ -z "$$GOPATH" ]; then \
		$(ECHO) "$(YELLOW)GOPATH is not set, defaulting to ~/go$(NC)"; \
		GOPATH=$$HOME/go; \
	fi
	@if ! echo "$$PATH" | grep -q "$$GOPATH/bin"; then \
		$(ECHO) "$(YELLOW)Warning: GOPATH/bin not found in PATH. Please add to your shell config:$(NC)"; \
		$(ECHO) "  export GOPATH=$$GOPATH"; \
		$(ECHO) "  export PATH=\$$PATH:\$$GOPATH/bin"; \
	fi

# API代码生成
api: $(API_TARGETS)

$(API_DEST_DIR)/%/. : $(API_SRC_DIR)/%.api
	@$(ECHO) "$(YELLOW)Generating API code for $<...$(NC)"
	@$(MKDIR) $(API_DEST_DIR)/$(notdir $(basename $<))
	@if $(GOZERO) api go -api $< -dir $(API_DEST_DIR)/$(notdir $(basename $<)); then \
		$(TOUCH) $@; \
		$(ECHO) "$(GREEN)API code generated to $(API_DEST_DIR)/$(notdir $(basename $<))$(NC)"; \
	else \
		$(ECHO) "$(RED)Failed to generate API code for $<$(NC)"; \
		exit 1; \
	fi

# RPC代码生成
rpc: $(RPC_TARGETS)

$(RPC_DEST_DIR)/%/. : $(RPC_SRC_DIR)/%.proto
	@$(ECHO) "$(YELLOW)Generating RPC code for $<...$(NC)"
	@$(MKDIR) $(RPC_DEST_DIR)/$(notdir $(basename $<))
	@if cd $(RPC_SRC_DIR) && \
	 $(GOZERO) rpc protoc $(notdir $<) \
  	 --proto_path=. \
  	 --proto_path=../../third_party \
	 --go_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --go-grpc_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --zrpc_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --plugin=protoc-gen-validate=$(GOPATH)/bin/protoc-gen-validate \
	 --style=go_zero \
	 -m; then \
		$(ECHO) "$(GREEN)RPC code generated to $(RPC_DEST_DIR)/$(notdir $(basename $<))$(NC)"; \
	else \
		$(ECHO) "$(RED)Failed to generate RPC code for $<$(NC)"; \
		exit 1; \
	fi
	@echo "Injecting custom tags to remove omitempty..."
	@find $(RPC_DEST_DIR) -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;

# 检查文件
check:
	@$(ECHO) "$(YELLOW)API files found:$(NC) $(API_FILES)"
	@$(ECHO) "$(YELLOW)Proto files found:$(NC) $(PROTO_FILES)"

# 帮助信息
help:
	@$(ECHO) "$(GREEN)Available targets:$(NC)"
	@$(ECHO) "  $(YELLOW)all$(NC)        - Generate all API and RPC code (default)"
	@$(ECHO) "  $(YELLOW)api$(NC)        - Generate API code"
	@$(ECHO) "  $(YELLOW)rpc$(NC)        - Generate RPC code"
	@$(ECHO) "  $(YELLOW)check$(NC)      - Check for API and Proto files"
	@$(ECHO) "  $(YELLOW)init$(NC)       - Initialize project environment (install tools, create dirs)"
	@$(ECHO) "  $(YELLOW)help$(NC)       - Show this help message"
	@$(ECHO) ""
	@$(ECHO) "$(GREEN)Examples:$(NC)"
	@$(ECHO) "  make init           # Initialize project environment"
	@$(ECHO) "  make                # Generate all code"
	@$(ECHO) "  make rpc            # Generate only RPC code"
	@$(ECHO) "  make check          # Check for API and Proto files"
