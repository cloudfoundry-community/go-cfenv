package cfenv

import (
	reporter "github.com/joefitzgerald/rainbow-reporter"
	"github.com/sclevine/spec"

	"testing"
)

var suite spec.Suite

func init() {
	suite = spec.New("go-cfenv internals", spec.Report(reporter.Rainbow{}))
	suite("envmap", testEnvMap)
}

func TestInternalsSuite(t *testing.T) {
	suite.Run(t)
}
