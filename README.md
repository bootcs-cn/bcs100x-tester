# BCS100X Tester

BCS100X 课程自动化测试工具。

## 使用

```bash
# 克隆项目
git clone https://github.com/bootcs-dev/bcs100x-tester
git clone https://github.com/bootcs-dev/bootcs-tester-utils

# 编译
cd bcs100x-tester
go build .

# 测试
./bcs100x-tester -s hello -d ~/cs50/hello
./bcs100x-tester --list
```

**环境依赖：** Go 1.24+, clang, python3, sqlite3

## License

MIT
