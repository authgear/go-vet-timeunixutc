package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/authgear/go-vet-timeunixutc/pkg/lint"
)

func main() {
	singlechecker.Main(lint.Analyzer)
}
