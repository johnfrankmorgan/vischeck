package main

import (
	"github.com/johnfrankmorgan/vischeck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(vischeck.Analyzer)
}
