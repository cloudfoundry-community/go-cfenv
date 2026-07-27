package cfenv

import (
	"github.com/sclevine/spec"

	"testing"
)

var suite spec.Suite

func init() {
	suite = spec.New("go-cfenv internals")
	suite("envmap", testEnvMap)
}

func TestInternalsSuite(t *testing.T) {
	suite.Run(t)
}
