package stages

import (
	"fmt"
	"time"

	"github.com/bootcs-dev/bcs100x-tester/internal/helpers"
	"github.com/bootcs-dev/tester-utils/runner"
	"github.com/bootcs-dev/tester-utils/test_case_harness"
	"github.com/bootcs-dev/tester-utils/tester_definition"
)

func marioLessTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "mario-less",
		Timeout:  30 * time.Second,
		TestFunc: testMarioLess,
	}
}

func testMarioLess(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir

	// 1. 检查 mario.c 文件存在
	logger.Infof("Checking mario.c exists...")
	if !harness.FileExists("mario.c") {
		return fmt.Errorf("mario.c does not exist")
	}
	logger.Successf("mario.c exists")

	// 2. 编译 mario.c
	logger.Infof("Compiling mario.c...")
	if err := helpers.CompileC(workDir, "mario.c", "mario", true); err != nil {
		return fmt.Errorf("mario.c does not compile: %v", err)
	}
	logger.Successf("mario.c compiles")

	// 3. 测试拒绝无效输入 (对齐 CS50 check50)
	rejectTests := []struct {
		input string
		name  string
	}{
		{"-1", "rejects a height of -1"},
		{"0", "rejects a height of 0"},
		{"foo", "rejects a non-numeric height of \"foo\""},
		{"", "rejects a non-numeric height of \"\""},
	}

	for _, tc := range rejectTests {
		logger.Infof("Testing %s...", tc.name)

		r := runner.Run(workDir, "mario").
			WithTimeout(5 * time.Second).
			Stdin(tc.input).
			Reject()

		if err := r.Error(); err != nil {
			return fmt.Errorf("%s: %v", tc.name, err)
		}

		logger.Successf("✓ %s", tc.name)
	}

	// 4. 测试有效输入
	validTests := []struct {
		height int
		name   string
	}{
		{1, "handles a height of 1 correctly"},
		{2, "handles a height of 2 correctly"},
		{8, "handles a height of 8 correctly"},
	}

	for _, tc := range validTests {
		logger.Infof("Testing %s...", tc.name)

		expected := helpers.GeneratePyramid(tc.height)
		input := fmt.Sprintf("%d", tc.height)

		r := runner.Run(workDir, "mario").
			WithTimeout(5 * time.Second).
			Stdin(input).
			Stdout(expected).
			Exit(0)

		if err := r.Error(); err != nil {
			return fmt.Errorf("%s: %v", tc.name, err)
		}

		logger.Successf("✓ %s", tc.name)
	}

	// 5. 测试拒绝后接受 (CS50 特有测试)
	logger.Infof("Testing rejects -1, then accepts 2...")
	expected := helpers.GeneratePyramid(2)
	r := runner.Run(workDir, "mario").
		WithTimeout(5 * time.Second).
		Stdin("-1\n2\n").
		Stdout(expected).
		Exit(0)

	if err := r.Error(); err != nil {
		return fmt.Errorf("rejects -1 then accepts 2: %v", err)
	}
	logger.Successf("✓ rejects -1, then accepts 2")

	logger.Successf("All mario-less tests passed!")
	return nil
}
