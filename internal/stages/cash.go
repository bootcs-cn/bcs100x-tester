package stages

import (
	"fmt"
	"time"

	"github.com/bootcs-dev/bcs100x-tester/internal/helpers"
	"github.com/bootcs-dev/tester-utils/runner"
	"github.com/bootcs-dev/tester-utils/test_case_harness"
	"github.com/bootcs-dev/tester-utils/tester_definition"
)

func cashTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "cash",
		Timeout:  30 * time.Second,
		TestFunc: testCash,
	}
}

func testCash(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir

	// 1. 检查 cash.c 文件存在
	logger.Infof("Checking cash.c exists...")
	if !harness.FileExists("cash.c") {
		return fmt.Errorf("cash.c does not exist")
	}
	logger.Successf("cash.c exists")

	// 2. 编译 cash.c
	logger.Infof("Compiling cash.c...")
	if err := helpers.CompileC(workDir, "cash.c", "cash", true); err != nil {
		return fmt.Errorf("cash.c does not compile: %v", err)
	}
	logger.Successf("cash.c compiles")

	// 3. 测试有效输入
	// 对齐 CS50 check50 的测试用例
	validTests := []struct {
		input    string
		expected string
		name     string
	}{
		{"41", "4", "input of 41 yields output of 4"},
		{"1", "1", "input of 1 yields output of 1"},
		{"15", "2", "input of 15 yields output of 2"},
		{"160", "7", "input of 160 yields output of 7"},
		{"2300", "92", "input of 2300 yields output of 92"},
	}

	for _, tc := range validTests {
		logger.Infof("Testing %s...", tc.name)

		r := runner.Run(workDir, "cash").
			WithTimeout(5 * time.Second).
			Stdin(tc.input).
			Stdout(tc.expected).
			Exit(0)

		if err := r.Error(); err != nil {
			return fmt.Errorf("%s: %v", tc.name, err)
		}

		logger.Successf("✓ %s", tc.name)
	}

	// 4. 测试拒绝无效输入 (对齐 CS50 check50)
	rejectTests := []struct {
		input string
		name  string
	}{
		{"-1", "rejects a negative input like -1"},
		{"foo", "rejects a non-numeric input of \"foo\""},
		{"", "rejects a non-numeric input of \"\""},
	}

	for _, tc := range rejectTests {
		logger.Infof("Testing %s...", tc.name)

		r := runner.Run(workDir, "cash").
			WithTimeout(5 * time.Second).
			Stdin(tc.input).
			Reject()

		if err := r.Error(); err != nil {
			return fmt.Errorf("%s: %v", tc.name, err)
		}

		logger.Successf("✓ %s", tc.name)
	}

	logger.Successf("All cash tests passed!")
	return nil
}
