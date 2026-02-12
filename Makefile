.PHONY: help test test-all lint format clean

help: ## 显示帮助信息
	@echo "Go Learning Path 项目构建工具"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## 运行当前目录的测试
	go test -v ./...

test-all: ## 运行所有测试
	go test -v ./...

lint: ## 代码检查
	golangci-lint run

format: ## 格式化代码
	go fmt ./...
	goimports -w .

clean: ## 清理构建文件
	go clean
	rm -rf bin/ cover.out

week1-day1: ## 运行第1周第1天的练习
	cd week1/day1 && go run goroutine_exercise.go

week1-day1-test: ## 测试第1周第1天的代码
	cd week1/day1 && go test -v

week1-day1-bench: ## 基准测试第1周第1天的代码
	cd week1/day1 && go test -bench=. -benchmem

# 开发工具检查
check-tools: ## 检查必要的开发工具
	@echo "检查Go版本:"
	@go version
	@echo ""
	@echo "检查golangci-lint:"
	@golangci-lint --version 2>/dev/null || echo "golangci-lint未安装，运行: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	@echo ""
	@echo "检查goimports:"
	@goimports --version 2>/dev/null || echo "goimports未安装，运行: go install golang.org/x/tools/cmd/goimports@latest"