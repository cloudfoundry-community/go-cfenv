package cfenv_test

import (
	reporter "github.com/joefitzgerald/rainbow-reporter"
	"github.com/sclevine/spec"

	"testing"
)

var suite spec.Suite

func init() {
	suite = spec.New("go-cfenv api", spec.Report(reporter.Rainbow{}))
	suite("environment", testEnvironment)
	suite("cfenv", testcfenv)
}

func TestSuite(t *testing.T) {
	suite.Run(t)
}
