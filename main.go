package main

import (
	"os"

	"github.com/bootcs-dev/bcs100x-tester/internal/stages"
	tester_utils "github.com/bootcs-dev/tester-utils"
)

func main() {
	definition := stages.GetDefinition()
	os.Exit(tester_utils.Run(os.Args[1:], definition))
}
