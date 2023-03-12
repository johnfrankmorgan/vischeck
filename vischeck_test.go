package vischeck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/johnfrankmorgan/vischeck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	wd, _ := os.Getwd()
	analysistest.Run(t, filepath.Join(wd, "testdata"), vischeck.Analyzer)
}
