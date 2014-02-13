# Go Cloud Foundry Environment Library (cfenv)

### Overview

cfenv is a library to assist you in writing Go apps that run on [Cloud Foundry](http://cloudfoundry.org). It provides convenience functions and structures that map to Cloud Foundry environment variable primitives (http://docs.cloudfoundry.com/docs/using/deploying-apps/environment-variable.html).

### Build Status

* [![Build Status - Master](https://travis-ci.org/joefitzgerald/cfenv.png?branch=master)](https://travis-ci.org/joefitzgerald/cfenv) `Master`
* [![Build Status - Develop](https://travis-ci.org/joefitzgerald/cfenv.png?branch=develop)](https://travis-ci.org/joefitzgerald/cfenv) `Develop`
* [![Build Status - Develop](https://travis-ci.org/joefitzgerald/cfenv.png?branch=develop)](https://travis-ci.org/joefitzgerald/cfenv)

### Usage

`go
package main

import (
	"github.com/joefitzgerald/cfenv"
)

func main() {
	appEnv := cfenv.Current()
	fmt.Printf("appEnv.Id")
	fmt.Printf("appEnv.Index")
	fmt.Printf("appEnv.Name")
	fmt.Printf("appEnv.Host")
	fmt.Printf("appEnv.Port")
	fmt.Printf("appEnv.Version")
	fmt.Printf("appEnv.Home")
	fmt.Printf("appEnv.MemoryLimit")
	fmt.Printf("appEnv.WorkingDir")
	fmt.Printf("appEnv.TempDir")
	fmt.Printf("appEnv.User")
	fmt.Printf("appEnv.Services")
}
`