.PHONY: all api rpc list-api list-rpc help clean

# Project directory configuration
API_DIR := cmd/api
PB_DIR := cmd/pb
APP_DIR := app

all: help

## Generate ALL API services
api:
	@echo Generating ALL API services...
	@if exist "$(API_DIR)\*.api" ( \
		for %%f in ("$(API_DIR)\*.api") do ( \
			echo Generating API service: %%~nf && \
			if not exist "$(APP_DIR)\%%~nf" mkdir "$(APP_DIR)\%%~nf" && \
			goctl api go -api "%%f" -dir "$(APP_DIR)\%%~nf" && \
			echo Success: Generated in $(APP_DIR)\%%~nf && \
			echo. \
		) \
	) else ( \
		echo No .api files found in $(API_DIR)/ \
	)

## Generate ALL RPC services
rpc:
	@echo Generating ALL RPC services...
	@if exist "$(PB_DIR)\*.proto" ( \
		for %%f in ("$(PB_DIR)\*.proto") do ( \
			echo Generating RPC service: %%~nf && \
			if not exist "$(APP_DIR)\%%~nf" mkdir "$(APP_DIR)\%%~nf" && \
			pushd "$(PB_DIR)" && \
			goctl rpc protoc "%%~nxf" --go_out="..\..\$(APP_DIR)\%%~nf" --go-grpc_out="..\..\$(APP_DIR)\%%~nf" --zrpc_out="..\..\$(APP_DIR)\%%~nf" -m && \
			popd && \
			echo Success: Generated in $(APP_DIR)\%%~nf && \
			echo. \
		) \
	) else ( \
		echo No .proto files found in $(PB_DIR)/ \
	)

## Generate specific API service
api-%:
	@echo Generating API service: $*
	@if not exist "$(APP_DIR)\$*" mkdir "$(APP_DIR)\$*"
	@if exist "$(API_DIR)\$*.api" ( \
		goctl api go -api "$(API_DIR)\$*.api" -dir "$(APP_DIR)\$*" && \
		echo Success: API service generated in $(APP_DIR)\$* \
	) else ( \
		echo Error: API definition file not found - $(API_DIR)\$*.api && \
		echo Available API services: && \
		$(MAKE) --no-print-directory list-api && \
		exit 1 \
	)

## Generate specific RPC service
rpc-%:
	@echo Generating RPC service: $*
	@if not exist "$(APP_DIR)\$*" mkdir "$(APP_DIR)\$*"
	@if exist "$(PB_DIR)\$*.proto" ( \
		pushd "$(PB_DIR)" && \
		goctl rpc protoc "$*.proto" --go_out="..\..\$(APP_DIR)\$*" --go-grpc_out="..\..\$(APP_DIR)\$*" --zrpc_out="..\..\$(APP_DIR)\$*" -m && \
		popd && \
		echo Success: RPC service generated in $(APP_DIR)\$* \
	) else ( \
		echo Error: Proto file not found - $(PB_DIR)\$*.proto && \
		echo Available RPC services: && \
		$(MAKE) --no-print-directory list-rpc && \
		exit 1 \
	)

## List available API services
list-api:
	@echo Available API services (from $(API_DIR)/):
	@if exist "$(API_DIR)\*.api" ( \
		for %%f in ("$(API_DIR)\*.api") do ( \
			echo   - %%~nf \
		) \
	) else ( \
		echo   No .api files found in $(API_DIR)/ \
	)

## List available RPC services
list-rpc:
	@echo Available RPC services (from $(PB_DIR)/):
	@if exist "$(PB_DIR)\*.proto" ( \
		for %%f in ("$(PB_DIR)\*.proto") do ( \
			echo   - %%~nf \
		) \
	) else ( \
		echo   No .proto files found in $(PB_DIR)/ \
	)

## Clean all generated files
clean:
	@echo Cleaning all generated services...
	@if exist "$(APP_DIR)" ( \
		for /d %%d in ("$(APP_DIR)\*") do ( \
			echo Deleting service: %%~nxd && \
			rmdir /s /q "%%d" \
		) \
	)
	@echo Clean complete.

## Show help information
help:
	@echo Go-Zero Project Generator
	@echo.
	@echo Usage:
	@echo   make api               Generate ALL API services
	@echo   make rpc               Generate ALL RPC services
	@echo   make api-^<name^>      Generate specific API service
	@echo   make rpc-^<name^>      Generate specific RPC service
	@echo   make list-api          List all API services
	@echo   make list-rpc          List all RPC services
	@echo   make clean             Remove all generated services
	@echo.
	@echo Examples:
	@echo   make api               Generate all API services
	@echo   make rpc-users         Generate specific RPC service
	@echo   make clean             Remove all generated code
	@echo.
	@echo Directory Structure:
	@echo   API definitions: cmd/api/*.api
	@echo   Proto files:     cmd/pb/*.proto
	@echo   Generated code:  app/^<service_name^>/

%::
	@rem