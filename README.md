# Go Cloud Foundry Environment Library (cfenv)

### Overview 

[![GoDoc](https://godoc.org/github.com/joefitzgerald/cfenv?status.png)](https://godoc.org/github.com/joefitzgerald/cfenv)

`cfenv` is a library to assist you in writing Go apps that run on [Cloud Foundry](http://cloudfoundry.org). It provides convenience functions and structures that map to Cloud Foundry environment variable primitives (http://docs.cloudfoundry.com/docs/using/deploying-apps/environment-variable.html).

### Build Status

* [![Build Status - Master](https://travis-ci.org/joefitzgerald/cfenv.png?branch=master)](https://travis-ci.org/joefitzgerald/cfenv) `Master`
* [![Build Status - Develop](https://travis-ci.org/joefitzgerald/cfenv.png?branch=develop)](https://travis-ci.org/joefitzgerald/cfenv) `Develop`

### Usage

`go get github.com/joefitzgerald/cfenv`

```go
package main

import (
	"github.com/joefitzgerald/cfenv"
)

func main() {
	appEnv := cfenv.Current()
	
	fmt.Println("ID:", appEnv.Id)
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