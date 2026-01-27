package stages

import (
	"github.com/bootcs-dev/tester-utils/tester_definition"
)

// GetDefinition 返回 tester 的完整定义
func GetDefinition() tester_definition.TesterDefinition {
	return tester_definition.TesterDefinition{
		TestCases: []tester_definition.TestCase{
			// Week 1: C 基础
			helloTestCase(),
			marioLessTestCase(),
			marioMoreTestCase(),
			cashTestCase(),
			creditTestCase(),
			// TODO: 添加更多测试用例
		},
	}
}
