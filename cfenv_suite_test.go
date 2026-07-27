package cfenv_test

import (
	"github.com/sclevine/spec"

	"testing"
)

var suite spec.Suite

func init() {
	suite = spec.New("go-cfenv api")
	suite("environment", testEnvironment)
	suite("cfenv", testcfenv)
}

func TestSuite(t *testing.T) {
	suite.Run(t)
}
