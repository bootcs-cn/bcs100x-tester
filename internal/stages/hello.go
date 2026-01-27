package stages

import (
	"fmt"
	"time"

	"github.com/bootcs-dev/bcs100x-tester/internal/helpers"
	"github.com/bootcs-dev/tester-utils/runner"
	"github.com/bootcs-dev/tester-utils/test_case_harness"
	"github.com/bootcs-dev/tester-utils/tester_definition"
)

func helloTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "hello",
		Timeout:  30 * time.Second,
		TestFunc: testHello,
	}
}

func testHello(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir

	// 1. 检查 hello.c 文件存在
	logger.Infof("Checking hello.c exists...")
	if !harness.FileExists("hello.c") {
		return fmt.Errorf("hello.c does not exist")
	}
	logger.Successf("hello.c exists")

	// 2. 编译 hello.c
	logger.Infof("Compiling hello.c...")
	if err := helpers.CompileC(workDir, "hello.c", "hello", true); err != nil {
		return fmt.Errorf("hello.c does not compile: %v", err)
	}
	logger.Successf("hello.c compiles")

	// 3. 测试用例：对齐 CS50 check50 官方测试
	testCases := []struct {
		name     string
		expected string
	}{
		{"Emma", "Emma"},
		{"Rodrigo", "Rodrigo"},
	}

	for _, tc := range testCases {
		logger.Infof("Testing with input %q...", tc.name)

		r := runner.Run(workDir, "hello").
			WithTimeout(5 * time.Second).
			Stdin(tc.name).
			Stdout(tc.expected).
			Exit(0)

		if err := r.Error(); err != nil {
			return fmt.Errorf("test failed for input %q: %v", tc.name, err)
		}

		logger.Successf("✓ Output correct for input %q", tc.name)
	}

	logger.Successf("All tests passed!")
	return nil
}
