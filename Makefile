# Go-Vue-Admin Makefile

.PHONY: build run swagger swag clean test

# 构建项目
build:
	go build -o admin-server main.go

# 运行项目
run:
	go run main.go

# 安装 swagger 工具
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# 生成 swagger 文档
swagger:
	swag init -g main.go -o ./docs

# 快捷命令别名
swag: swagger

# 初始化数据库
db:
	go run main.go -db

# 重置数据库
reset-db:
	go run main.go -reset-db

# 重置管理员密码
reset-pwd:
	go run main.go -reset-pwd

# 清理编译文件
clean:
	rm -f admin-server
	rm -rf ./docs/swagger.json ./docs/swagger.yaml

# 运行测试
test:
	go test -v ./...

# 下载依赖
deps:
	go mod tidy
	go mod download

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  make build       - 构建项目"
	@echo "  make run         - 运行项目"
	@echo "  make swagger     - 生成 swagger 文档"
	@echo "  make install-swag - 安装 swagger 工具"
	@echo "  make db          - 初始化数据库"
	@echo "  make reset-db    - 重置数据库"
	@echo "  make reset-pwd   - 重置管理员密码"
	@echo "  make clean       - 清理编译文件"
	@echo "  make test        - 运行测试"
	@echo "  make deps        - 下载依赖"
