package stages

import (
	"fmt"
	"time"

	"github.com/bootcs-cn/bcs100x-tester/internal/helpers"
	"github.com/bootcs-cn/tester-utils/runner"
	"github.com/bootcs-cn/tester-utils/test_case_harness"
	"github.com/bootcs-cn/tester-utils/tester_definition"
)

func caesarTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "caesar",
		Timeout:  30 * time.Second,
		TestFunc: testCaesar,
	}
}

func testCaesar(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir

	// 1. 检查 caesar.c 文件存在
	logger.Infof("Checking caesar.c exists...")
	if !harness.FileExists("caesar.c") {
		return fmt.Errorf("caesar.c does not exist")
	}
	logger.Successf("caesar.c exists")

	// 2. 编译 caesar.c
	logger.Infof("Compiling caesar.c...")
	if err := helpers.CompileC(workDir, "caesar.c", "caesar", true); err != nil {
		return fmt.Errorf("caesar.c does not compile: %v", err)
	}
	logger.Successf("caesar.c compiles")

	// 3. 加密测试用例（完全对齐 CS50 check50）
	encryptTests := []struct {
		key        string
		plaintext  string
		ciphertext string
		name       string
	}{
		{
			"1", "a", "b",
			"encrypts 'a' as 'b' using 1 as key",
		},
		{
			"23", "barfoo", "yxocll",
			"encrypts 'barfoo' as 'yxocll' using 23 as key",
		},
		{
			"3", "BARFOO", "EDUIRR",
			"encrypts 'BARFOO' as 'EDUIRR' using 3 as key",
		},
		{
			"4", "BaRFoo", "FeVJss",
			"encrypts 'BaRFoo' as 'FeVJss' using 4 as key",
		},
		{
			"65", "barfoo", "onesbb",
			"encrypts 'barfoo' as 'onesbb' using 65 as key",
		},
		{
			"12", "world, say hello!", "iadxp, emk tqxxa!",
			"encrypts 'world, say hello!' as 'iadxp, emk tqxxa!' using 12 as key",
		},
	}

	for _, tc := range encryptTests {
		logger.Infof("Testing %s...", tc.name)

		r := runner.Run(workDir, "caesar", tc.key).
			WithTimeout(5 * time.Second).
			Stdin(tc.plaintext).
			Stdout(tc.ciphertext).
			Exit(0)

		if err := r.Error(); err != nil {
			return fmt.Errorf("test failed for %s: %v", tc.name, err)
		}

		logger.Successf("✓ %s", tc.name)
	}

	// 4. 错误处理测试用例
	errorTests := []struct {
		args []string
		name string
	}{
		{
			[]string{},
			"handles lack of argv[1]",
		},
		{
			[]string{"2x"},
			"handles non-numeric key",
		},
		{
			[]string{"1", "2"},
			"handles too many arguments",
		},
	}

	for _, tc := range errorTests {
		logger.Infof("Testing %s...", tc.name)

		r := runner.Run(workDir, "caesar", tc.args...).
			WithTimeout(5 * time.Second).
			Execute().
			Exit(1)

		if err := r.Error(); err != nil {
			return fmt.Errorf("test failed for %s: %v", tc.name, err)
		}

		logger.Successf("✓ %s", tc.name)
	}

	logger.Successf("All caesar tests passed!")
	return nil
}
