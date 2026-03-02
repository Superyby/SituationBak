.PHONY: all build clean proto test run-gateway run-auth run-user run-satellite run-favorite docker-up docker-down

# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Proto相关
PROTOC=protoc
PROTO_DIR=api/proto
PROTO_GEN_DIR=api/gen

# 服务目录
SERVICES=gateway auth user satellite favorite

# 构建输出目录
BUILD_DIR=bin

all: build

# 构建所有服务
build: $(SERVICES)

gateway:
	$(GOBUILD) -o $(BUILD_DIR)/gateway ./services/gateway/cmd

auth:
	$(GOBUILD) -o $(BUILD_DIR)/auth-svc ./services/auth/cmd

user:
	$(GOBUILD) -o $(BUILD_DIR)/user-svc ./services/user/cmd

satellite:
	$(GOBUILD) -o $(BUILD_DIR)/satellite-svc ./services/satellite/cmd

favorite:
	$(GOBUILD) -o $(BUILD_DIR)/favorite-svc ./services/favorite/cmd

# 清理构建产物
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(PROTO_GEN_DIR)

# 生成Proto代码
proto:
	@echo "Generating protobuf code..."
	@mkdir -p $(PROTO_GEN_DIR)
	$(PROTOC) --go_out=$(PROTO_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) --go-grpc_opt=paths=source_relative \
		-I$(PROTO_DIR) \
		$(PROTO_DIR)/auth/v1/*.proto \
		$(PROTO_DIR)/user/v1/*.proto \
		$(PROTO_DIR)/satellite/v1/*.proto \
		$(PROTO_DIR)/favorite/v1/*.proto

# 下载依赖
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 运行测试
test:
	$(GOTEST) -v ./...

# 本地运行服务
run-gateway:
	$(GOCMD) run ./services/gateway/cmd -config=./services/gateway/configs/config.yaml

run-auth:
	$(GOCMD) run ./services/auth/cmd -config=./services/auth/configs/config.yaml

run-user:
	$(GOCMD) run ./services/user/cmd -config=./services/user/configs/config.yaml

run-satellite:
	$(GOCMD) run ./services/satellite/cmd -config=./services/satellite/configs/config.yaml

run-favorite:
	$(GOCMD) run ./services/favorite/cmd -config=./services/favorite/configs/config.yaml

# Docker相关命令
docker-up:
	docker-compose -f deployments/docker/docker-compose.yaml up -d

docker-down:
	docker-compose -f deployments/docker/docker-compose.yaml down

docker-logs:
	docker-compose -f deployments/docker/docker-compose.yaml logs -f

docker-build:
	docker-compose -f deployments/docker/docker-compose.yaml build

# 启动基础设施
infra-up:
	docker-compose -f deployments/docker/docker-compose.yaml up -d mysql redis consul

infra-down:
	docker-compose -f deployments/docker/docker-compose.yaml down mysql redis consul

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all          - Build all services (default)"
	@echo "  build        - Build all services"
	@echo "  clean        - Clean build artifacts"
	@echo "  proto        - Generate protobuf code"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  test         - Run tests"
	@echo "  run-gateway  - Run gateway service locally"
	@echo "  run-auth     - Run auth service locally"
	@echo "  run-user     - Run user service locally"
	@echo "  run-satellite- Run satellite service locally"
	@echo "  run-favorite - Run favorite service locally"
	@echo "  docker-up    - Start all services with Docker"
	@echo "  docker-down  - Stop all Docker services"
	@echo "  docker-logs  - View Docker logs"
	@echo "  infra-up     - Start infrastructure only (MySQL, Redis, Consul)"
	@echo "  infra-down   - Stop infrastructure"
