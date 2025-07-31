# 工具检查
GO := go
GOZERO := goctl
MKDIR := mkdir -p
RM := rm -rf
TOUCH := touch

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

.PHONY: all api rpc clean force

all: api rpc

# API代码生成 - 使用标记文件确保每次都会执行
api: $(API_TARGETS)

$(API_DEST_DIR)/%/. : $(API_SRC_DIR)/%.api
	@echo "Generating API code for $<..."
	@$(MKDIR) $(API_DEST_DIR)/$(notdir $(basename $<))
	@$(GOZERO) api go -api $< -dir $(API_DEST_DIR)/$(notdir $(basename $<))
	@$(TOUCH) $(API_DEST_DIR)/$(notdir $(basename $<))/.
	@echo "API code generated to $(API_DEST_DIR)/$(notdir $(basename $<))"

# RPC代码生成 - 使用标记文件确保每次都会执行
rpc: $(RPC_TARGETS)

$(RPC_DEST_DIR)/%/. : $(RPC_SRC_DIR)/%.proto force
	@echo "Generating RPC code for $<..."
	@$(MKDIR) $(RPC_DEST_DIR)/$(notdir $(basename $<))
	@cd $(RPC_SRC_DIR) && \
	 $(GOZERO) rpc protoc $(notdir $<) \
	 --go_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --go-grpc_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 --zrpc_out=../../$(RPC_DEST_DIR)/$(notdir $(basename $<)) \
	 -m
	@$(TOUCH) $(RPC_DEST_DIR)/$(notdir $(basename $<))/.
	@echo "RPC code generated to $(RPC_DEST_DIR)/$(notdir $(basename $<))"

# 检查文件
check:
	@echo "API files found: $(API_FILES)"
	@echo "Proto files found: $(PROTO_FILES)"

# 清理生成代码
clean:
	@echo "Cleaning generated code..."
	@$(RM) $(API_DEST_DIR) $(RPC_DEST_DIR)
	@echo "Clean complete"