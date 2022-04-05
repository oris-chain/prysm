package uintcast_test

import (
	"testing"

	"github.com/prysmaticlabs/prysm/build/bazel"
	"github.com/prysmaticlabs/prysm/tools/analyzers/uintcast"
	"golang.org/x/tools/go/analysis/analysistest"
)

func init() {
	if bazel.BuiltWithBazel() {
		bazel.SetGoEnv()
	}
}

func TestAnalyzer(t *testing.T) {
	testdata := bazel.TestDataPath(t)
	analysistest.TestData = func() string { return testdata }
	analysistest.Run(t, testdata, uintcast.Analyzer)
}
