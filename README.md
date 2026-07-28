# Go CF Environment Package

[![CI](https://github.com/cloudfoundry-community/go-cfenv/actions/workflows/ci.yml/badge.svg)](https://github.com/cloudfoundry-community/go-cfenv/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/github/actions/workflow/status/cloudfoundry-community/go-cfenv/codeql.yml?label=CodeQL)](https://github.com/cloudfoundry-community/go-cfenv/security/code-scanning)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/cloudfoundry-community/go-cfenv/badge)](https://scorecard.dev/viewer/?uri=github.com/cloudfoundry-community/go-cfenv)
[![Go Version](https://img.shields.io/github/go-mod/go-version/cloudfoundry-community/go-cfenv)](go.mod)
[![Latest Release](https://img.shields.io/github/v/release/cloudfoundry-community/go-cfenv)](https://github.com/cloudfoundry-community/go-cfenv/releases)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

### Overview

`cfenv` is a package to assist you in writing Go apps that run on [Cloud Foundry](http://cloudfoundry.org). It provides convenience functions and structures that map to Cloud Foundry environment variable primitives (http://docs.cloudfoundry.org/devguide/deploy-apps/environment-variable.html).

### Requirements

Go 1.18 or later.

### Usage

`go get github.com/cloudfoundry-community/go-cfenv`

```go
package main

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfenv"
)

func main() {
	appEnv, _ := cfenv.Current()

	fmt.Println("ID:", appEnv.ID)
	fmt.Println("Index:", appEnv.Index)
	fmt.Println("Name:", appEnv.Name)
	fmt.Println("Host:", appEnv.Host)
	fmt.Println("Port:", appEnv.Port)
	fmt.Println("Version:", appEnv.Version)
	fmt.Println("Home:", appEnv.Home)
	fmt.Println("MemoryLimit:", appEnv.MemoryLimit)
	fmt.Println("WorkingDir:", appEnv.WorkingDir)
	fmt.Println("TempDir:", appEnv.TempDir)
	fmt.Println("User:", appEnv.User)
	fmt.Println("Services:", appEnv.Services)
}
```

### Contributing

Pull requests welcomed. Please run `make check` and `make audit` before
submitting.
