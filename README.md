# Synopsis

Package log implements a simple structured JSON logger. The goal of this package is to eliminate logging clutter generated by more complex logging packages that encourage useless type conversions, declarations, and other noise. The logger starts with a few pre-initialized log fields that can be extended in the package scope, of which individual log lines can be derived from and operated on.

# Variables

To use, first override the package-scoped variables at runtime.

```
var (
	// Service name
	Service = ""

	// Time is your time function
	Time    = func() interface{} {
		return time.Now().UnixNano() / int64(time.Millisecond)
	}

	// Default is the level used when calling Printf
	Default = Info
)
```

# Install

```
go get github.com/as/log
go test github.com/as/log -v -bench . 
go test github.com/as/log -race -count 1
```

This code may also be copied and pasted into your microservice
and modified to your liking. Put it in a package called
log. A little copying is better than a little dependency.

The only file you need is `log.go`, there are no external dependencies.

# Example 1

main.go
```
package main

import "github.com/as/log"

func main() {
	log.Service = "ex"
	log.Time = func() interface{} { return "2121.12.04" }

	log.Error.Add(
		"env", "prod",
		"burning", true,
		"pi", 3.14,
	).Printf("error: %v", io.EOF)
}
```

output
```
{"svc":"ex", "time":"2121.12.04", "level":"error", "msg":"error: EOF", "env":"prod", "burning":true, "pi":3.14}
```

# Example 2

```
package main

import (
        "io"
        "github.com/as/log"
)

func main() {
        log.Service = "plumber"
        log.Tags = log.Tags.Add("env", "dev", "host", "scruffy") // package scoped tags attach to all logs
        unclog()
}

func unclog() {
        line := log.Info.Add(
                "op", "unclog",
                "account", "514423351351",
                "toilets", []string{"foo", "bar"}, // easily queried for membership in common log aggregators
        )
        line.Printf("enter")
        defer line.Printf("exit")

        // each operation copies the object and returns its copy with any tags added
        line.Add("action", "plunge").Printf("good god man")
        line.Error().Add("action", "flush", "err", io.ErrClosedPipe).Printf("damn")

        line = line.Warn().Add("action", "bill", "mood", "exhausted")
        line.Printf("you need a bigger toilet")
        line.Printf("it gets clogged to easily")
}
```

```
go mod init
go get github.com/as/log
touch main.go # copy the example above to this file
go run main.go
```

```
{"svc":"plumber", "ts":1682559034, "level":"info", "env":"dev", "host":"scruffy", "op":"unclog", "account":"514423351351", "toilets":["foo","bar"], "msg":"enter"}
{"svc":"plumber", "ts":1682559034, "level":"info", "env":"dev", "host":"scruffy", "op":"unclog", "account":"514423351351", "toilets":["foo","bar"], "action":"plunge", "msg":"good god man"}
{"svc":"plumber", "ts":1682559034, "level":"error", "env":"dev", "host":"scruffy", "op":"unclog", "account":"514423351351", "toilets":["foo","bar"], "action":"flush", "err":"io: read/write on closed pipe", "msg":"damn"}
{"svc":"plumber", "ts":1682559034, "level":"warn", "env":"dev", "host":"scruffy", "op":"unclog", "account":"514423351351", "toilets":["foo","bar"], "action":"bill", "mood":"exhausted", "msg":"you need a bigger toilet"}
{"svc":"plumber", "ts":1682559034, "level":"warn", "env":"dev", "host":"scruffy", "op":"unclog", "account":"514423351351", "toilets":["foo","bar"], "action":"bill", "mood":"exhausted", "msg":"it gets clogged to easily"}
{"svc":"plumber", "ts":1682559034, "level":"info", "env":"dev", "host":"scruffy", "op":"unclog", "account":"514423351351", "toilets":["foo","bar"], "msg":"exit"}
```

# Example 3: But my caller frames!

You can use `AddFunc`, which has behavior similar to `Add`. It returns a copy of the line with the attached function executing before every print operation. Some precautions have been taken to avoid infinite recursion, but use of this feature is still mildly discouraged.

```
package main

import (
        "fmt"
        "path"
        "runtime"

        "github.com/as/log"
)

func main() {
        log.Service = "plumber"
        unclog()
}

func unclog() {
        fn := func(l log.Line) log.Line {
                return l.Add("func", where())
        }
        line := log.Info.AddFunc(fn)
        line.Printf("enter")
        defer line.Printf("exit")
}

func where() string {
        pc, file, line, ok := runtime.Caller(1)
        if !ok {
                return "/dev/null"
        }
        file = path.Base(file)

        fn := runtime.FuncForPC(pc)
        if fn == nil {
                return fmt.Sprintf("%s:%d", file, line)
        }
        return fmt.Sprintf("%s:/%s/", file, fn.Name())
}
```
