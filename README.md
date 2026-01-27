# BCS100X Tester

基于 `bootcs-tester-utils` 的 BCS100X 课程自动测试器。

## 项目结构

```
bcs100x-tester/
├── main.go                     # 入口文件
├── go.mod                      # Go 模块定义
├── Makefile                    # 构建脚本
├── README.md                   # 项目说明
├── internal/
│   ├── stages/
│   │   ├── stages.go          # 所有阶段定义
│   │   ├── hello.go           # hello 阶段测试
│   │   └── ...                # 其他阶段测试
│   └── helpers/
│       ├── c_compiler.go      # C 编译辅助函数
│       └── ...                # 其他辅助函数
└── fixtures/                   # 测试 fixtures (可选)
    └── hello/
        ├── pass/              # 通过测试的样例
        └── fail/              # 失败测试的样例
```

## 开发

### 安装依赖

```bash
make deps
```

### 构建

```bash
make build
```

### 运行测试

```bash
# 在 bcs100x-starter/hello 目录中运行 hello 测试
cd ../bcs100x-starter/hello
../../bcs100x-tester/bcs100x-tester test hello

# 或者使用 go run
cd ../../bcs100x-tester
go run . test hello --dir ../bcs100x-starter/hello
```

### 代码格式化

```bash
make fmt
```

## 实现进度

- [x] hello - Week 1: C 基础
- [ ] mario-less - Week 1: C 基础
- [ ] mario-more - Week 1: C 基础
- [ ] cash - Week 1: C 基础
- [ ] credit - Week 1: C 基础
- [ ] caesar - Week 2: Arrays
- [ ] ... (其他 32 个阶段)

## 测试用例实现指南

每个 stage 的测试用例应该：

1. **检查文件存在** - 确保必需的源文件存在
2. **编译/准备** - 编译 C 代码或准备运行环境
3. **运行测试用例** - 使用不同的输入测试程序行为
4. **验证输出** - 检查程序输出是否符合预期
5. **验证退出码** - 确保程序正确退出

### C 语言题目模板

```go
func exampleTestCase() tester_definition.TestCase {
    return tester_definition.TestCase{
        Slug:     "example",
        Timeout:  30 * time.Second,
        TestFunc: testExample,
    }
}

func testExample(harness *test_case_harness.TestCaseHarness) error {
    logger := harness.Logger

    // 1. 检查文件存在
    logger.Infof("Checking example.c exists...")
    if !helpers.FileExists(harness, "example.c") {
        return fmt.Errorf("example.c does not exist")
    }

    // 2. 编译
    logger.Infof("Compiling example.c...")
    if err := helpers.CompileC(harness, "example.c", "example", true); err != nil {
        return fmt.Errorf("example.c does not compile: %v", err)
    }

    // 3. 运行测试用例
    testCases := []struct {
        input    string
        expected string
    }{
        {"input1", "output1"},
        {"input2", "output2"},
    }

    for _, tc := range testCases {
        executable := harness.NewExecutable()
        executable.Command = filepath.Join(harness.WorkingDir, "example")

        if err := executable.Run(); err != nil {
            return err
        }

        executable.SendStdin(tc.input + "\n")
        result, err := executable.Wait()

        if err != nil {
            return err
        }

        if result.ExitCode != 0 {
            return fmt.Errorf("expected exit code 0, got %d", result.ExitCode)
        }

        output := strings.TrimSpace(result.Stdout)
        if !strings.Contains(output, tc.expected) {
            return fmt.Errorf("expected output to contain %q, got %q", tc.expected, output)
        }
    }

    return nil
}
```

## 参考资料

- [bootcs-tester-utils](https://github.com/bootcs-dev/tester-utils)
- [CS50 2025/x Problems](https://cs50.harvard.edu/x/2025/)
- [check50](https://github.com/cs50/check50)
