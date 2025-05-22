# 示例
GOOS ?= linux
GOARCH ?= amd64

# proto 生成
proto:
	protoc --go_out=. --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		proto/*/*.proto

.PHONY: proto

# build
build:
	rm -rf ./bin
	mkdir -p bin/
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/ ./...
	@if command -v upx >/dev/null 2>&1; then \
		upx ./bin/*; \
	else \
		echo "upx 未安装，跳过压缩"; \
	fi

.PHONY: build

# 初始化依赖（ent相关）
init:
	go get entgo.io/ent/entc/gen@v0.14.1
	go get entgo.io/ent/cmd/internal/printer@v0.14.1
	go get entgo.io/ent/cmd/ent@v0.14.1

.PHONY: init

# generate ent code
ent:
	@go run entgo.io/ent/cmd/ent generate \
		--feature privacy \
		--feature sql/modifier \
		--feature intercept,schema/snapshot \
		--feature entql \
		--feature sql/upsert \
		./internal/data/ent/schema

.PHONY: ent

# 本地开发用 server
dev-server:
	@go run ./cmd/server/main.go

# 本地开发用 agent
dev-agent:
	@go run ./cmd/agent/main.go

# 快速本地 dev 提示
dev:
	@echo "打开两个终端，分别执行："
	@echo "make dev-server"
	@echo "make dev-agent"

.PHONY: dev dev-server dev-agent

# 清理
clean:
	rm -rf ./bin
	find ./proto -name "*.pb.go" -delete
.PHONY: clean

