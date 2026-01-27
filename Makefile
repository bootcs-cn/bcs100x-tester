.PHONY: build test clean run install uninstall

# 构建 tester
build:
	go build -o bcs100x-tester .

# 安装到 $GOPATH/bin (可全局使用)
install:
	go install .

# 卸载
uninstall:
	rm -f $(shell go env GOPATH)/bin/bcs100x-tester

# 运行测试
test:
	go test -v ./...

# 清理构建产物
clean:
	rm -f bcs100x-tester
	go clean

# 运行 tester (需要指定工作目录)
run:
	go run . $(ARGS)

# 安装依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...
