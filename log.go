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
//
package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	// Service name
	Service = ""

	// Level is the default log level
	Level   = "info"

	// Time is your time function
	Time    = func() string {
		return fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))
	}
)

// Warn warns
func Warn(f string, v ...interface{}) { line{level: "warn"}.Printf(f, v...) }

// Info informs
func Info(f string, v ...interface{}) { line{level: "info"}.Printf(f, v...) }

// Fatal also informs, then panics
func Fatal(f string, v ...interface{}) {
	line{level: "fatal"}.Printf(f, v...)
	panic("fatal error") // no os.Exit allowed; doesn't run defer funcs
}

// Line returns a log line from a list of fields. Fields should be
// provided in a pair of two, otherwise the remainder is ignored.
//
// Empty field values are also ignored.
//
// Example:
//
// Line(
// 	"id", 1,
// 	"name", "x",
// )
//
func Line(fields ...interface{}) (msg string) {
	// If you're worried about type saftey, please remember
	// it is a logger.
	sep := ""
	for i := 0; i+1 < len(fields); i += 2 {
		key, val := fields[i], fields[i+1]
		if val == "" || val == nil {
			continue
		}
		msg += fmt.Sprintf(`%s%q:%s`, sep, key, quote(val))
		sep = ", "
	}
	return "{" + msg + "}"
}

type line struct {
	f     string
	err   error
	level string
}

func (l line) Printf(f string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, l.Sprintf(f, v...))
}
func (l line) Sprintf(f string, v ...interface{}) string {
	if l.level == "" {
		l.level = Level
	}
	return Line(
		"svc", Service,
		"time", Time(),
		"level", l.level,
		"err", l.err,
		"msg", fmt.Sprintf(f, v...),
	)
}

func quote(v interface{}) string {
	if v == nil {
		v = ""
	}
	switch v.(type) {
	case fmt.Stringer, error:
		v = fmt.Sprint(v)
	}
	data, _ := json.Marshal(v)
	return string(data)
}
