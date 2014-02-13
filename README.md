# Go Cloud Foundry Environment Library (cfenv)

### Overview

cfenv is a library to assist you in writing Go apps that run on [Cloud Foundry](http://cloudfoundry.org). It provides convenience functions and structures that map to Cloud Foundry environment variable primitives (http://docs.cloudfoundry.com/docs/using/deploying-apps/environment-variable.html).

### Build Status

* [![Build Status - Master](https://travis-ci.org/joefitzgerald/cfenv.png?branch=master)](https://travis-ci.org/joefitzgerald/cfenv) `Master`
* [![Build Status - Develop](https://travis-ci.org/joefitzgerald/cfenv.png?branch=develop)](https://travis-ci.org/joefitzgerald/cfenv) `Develop`

### Usage

```go
package main

import (
	"github.com/joefitzgerald/cfenv"
)

func main() {
	appEnv := cfenv.Current()
	fmt.Printf("Id:", appEnv.Id)
	fmt.Printf("Index:", appEnv.Index)
	fmt.Printf("Name:", appEnv.Name)
	fmt.Printf("Host:", appEnv.Host)
	fmt.Printf("Port:", appEnv.Port)
	fmt.Printf("Version:", appEnv.Version)
	fmt.Printf("Home:", appEnv.Home)
	fmt.Printf("MemoryLimit:", appEnv.MemoryLimit)
	fmt.Printf("WorkingDir:", appEnv.WorkingDir)
	fmt.Printf("TempDir:", appEnv.TempDir)
	fmt.Printf("User:", appEnv.User)
	fmt.Printf("Services:", appEnv.Services)
}
```