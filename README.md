# BCS100X Tester

BCS100X 课程自动化测试工具。

## 方式一：源码运行

```bash
git clone https://github.com/bootcs-cn/bcs100x-tester
cd bcs100x-tester
go build .
./bcs100x-tester -s hello -d ~/my-solution/hello
```

**环境依赖：** Go 1.24+, clang, python3, sqlite3

## 方式二：Docker 镜像

**快速开始**

```bash
cd ~/my-solution  # 你的代码根目录
docker pull bootcs/bcs100x-tester
docker run --rm -v "$(pwd):/workspace" bootcs/bcs100x-tester -s hello -d /workspace/hello
```

**简化脚本（推荐）**

在代码根目录创建 `test.sh`：

```bash
#!/bin/bash
docker run --rm -v "$(pwd):/workspace" bootcs/bcs100x-tester \
  -s "${1:-hello}" -d "/workspace/${1:-hello}"
```

使用：`chmod +x test.sh && ./test.sh hello`

**本地构建（可选）**

```bash
git clone https://github.com/bootcs-cn/bcs100x-tester
cd bcs100x-tester
docker build -t my-tester .
# 使用: docker run --rm -v ~/my-solution:/workspace my-tester -s hello -d /workspace/hello
```

## License

MIT
