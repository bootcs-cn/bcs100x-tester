module github.com/bootcs-cn/bcs100x-tester

go 1.24

toolchain go1.24.12

require (
	github.com/bootcs-cn/tester-utils v1.0.0
	github.com/mattn/go-sqlite3 v1.14.33
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/creack/pty v1.1.24 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// 本地开发时使用：go mod edit -replace github.com/bootcs-cn/tester-utils=../../bootcs-tester-utils
// replace github.com/bootcs-cn/tester-utils => ../../bootcs-tester-utils
