package main

import (
	"github.com/johnfrankmorgan/vischeck"
	"golang.org/x/tools/go/analysis"
)

var AnalyzerPlugin = analyzerPlugin{}

type analyzerPlugin struct{}

func (plugin analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		vischeck.Analyzer,
	}
}
