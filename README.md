# log

```
package main

import "github.com/as/log"

func main(){
	log.Service = "test"
	log.Info("Hello, Playground")
	log.Warn("Hello, Playground")
	log.Fatal("Hello, Playground")
}

{"svc":"test", "time":1621481516021, "level":"info", "msg":"Hello, Playground"}
{"svc":"test", "time":1621481516021, "level":"warn", "msg":"Hello, Playground"}
{"svc":"test", "time":1621481516021, "level":"fatal", "msg":"Hello, Playground"}
panic: fatal error

goroutine 1 [running]:

```

# description

```
// Package log implements a structured JSON log which can't be
// dependency injected into your microservice.
//
// To use the log, override the package-scoped variables.
// No need to use dependency injection. Dependency injecting
// a log only to configure it per-process in a microservice is
// often a beauracratic, unnecessary practice.
//
// This code may be copied and pasted into your microservice
// and modified to your liking. Put it in a package called
// log. A little copying is better than a little dependency.
```

# variables

```
var (
	// Service name
	Service = ""

	// Level is the default log level
	Level   = "info"

	// Time is your time function
	Time    = func() interface{} {
		return time.Now().UnixNano() / int64(time.Millisecond)
	}
)
```
