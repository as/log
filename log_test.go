package log_test

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/as/log"
)

func init() {
	log.Service = "test"
	log.Time = func() interface{} { return 12345 }
}

func TestLog(t *testing.T) {
	have := log.Info.Msg("test log message: %s", "package scoped variables arent hard to test").String()
	want := `{"svc":"test", "ts":12345, "level":"info", "msg":"test log message: package scoped variables arent hard to test"}`
	if have != want {
		t.Fatalf("bad log:\n\t\thave: %s\n\t\twant: %s", have, want)
	}
}

func TestAdd(t *testing.T) {
	have := log.Error.Add(
		"ip", "1.2.3.4",
		"port", "1111",
		"client", "mothra",
		"host", "example.com",
		"path", "/file.txt",
		"query", "what",
		"err", io.EOF,
	).Msg("custom fields").String()
	want := `{"svc":"test", "ts":12345, "level":"error", "ip":"1.2.3.4", "port":"1111", "client":"mothra", "host":"example.com", "path":"/file.txt", "query":"what", "err":"EOF", "msg":"custom fields"}`
	if have != want {
		t.Fatalf("bad log:\n\t\thave: %s\n\t\twant: %s", have, want)
	}
}

func TestTag(t *testing.T) {
	before := log.Tags
	log.Tags = log.Tags.Add("subcmd", "test")
	defer func() {
		log.Tags = before
	}()

	have := log.Error.Add("ip", "1.2.3.4").Msg("custom tags").String()
	want := `{"svc":"test", "ts":12345, "level":"error", "subcmd":"test", "ip":"1.2.3.4", "msg":"custom tags"}`
	if have != want {
		t.Fatalf("bad log:\n\t\thave: %s\n\t\twant: %s", have, want)
	}
}

func TestRace(t *testing.T) {
	defer log.SetOutput(log.SetOutput(ioutil.Discard))

	wg := sync.WaitGroup{}
	defer wg.Wait()

	ln := log.Info.Add("test", "TestRace")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		ln := ln.Add("proc", i)
		go func() {
			for i := 0; i < 1000*100; i++ {
				ln.Printf("count: %d", i)
			}
			wg.Done()
		}()
	}
}

func TestFatal(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("didnt panic")
		}
	}()
	defer log.SetOutput(log.SetOutput(ioutil.Discard))
	log.Fatal.F("panic: %v", io.EOF)
}

func TestExport(t *testing.T) {
	before := log.Tags
	log.Tags = log.Tags.Add("env", "dev", "version", 1, "git", "af753", "empty", "", "", 6)
	defer func() {
		log.Tags = before
	}()
	want := "env,dev,version,1,git,af753"
	have := strings.Join(log.Tags.Export(), ",")
	if have != want {
		t.Fatalf("bad log:\n\t\thave: %s\n\t\twant: %s", have, want)
	}
}

func TestAddFunc(t *testing.T) {
	line := log.Info.Add("test", "TestAddFunc")
	ctr := 0
	fn := func(l log.Line) log.Line {
		ctr++
		t.Log(l.String())
		l.Printf("nested recursive call")
		// doing a line.Printf here would crash the program
		return l
	}
	line.String() // 0
	line.String() // 0

	line = line.Error().AddFunc(fn).Add("test", "TestAddFunc")
	line.String()              // 1
	line = line.Msg("still 1") // no op
	line.Printf("2")
	line.F("3")

	line2 := line.AddFunc(nil)
	line2.String()
	line2.String()

	line.F("4")
	if ctr != 4 {
		t.Fatalf("bad count, want 4 have %d", ctr)
	}
}

func TestAddArray(t *testing.T) {
	hint := []string{}
	hint = nil
	have := log.Error.Add(
		"ip", "1.2.3.4",
		"port", "1111",
		"client", "mothra",
		"host", "example.com",
		"path", "/file.txt",
		"query", "what",
		"err", io.EOF,
		"hint", []string{},
		"null", hint,
	).Msg("custom fields").String()
	want := `{"svc":"test", "ts":12345, "level":"error", "ip":"1.2.3.4", "port":"1111", "client":"mothra", "host":"example.com", "path":"/file.txt", "query":"what", "err":"EOF", "msg":"custom fields"}`
	if have != want {
		t.Fatalf("bad log:\n\t\thave: %s\n\t\twant: %s", have, want)
	}
}

func BenchmarkLog(b *testing.B) {
	defer log.SetOutput(log.SetOutput(ioutil.Discard))

	b.Run("Printf", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			log.Printf("count: %d", b.N)
		}
	})

	b.Run("Print20", func(b *testing.B) {
		ln := log.Error.Add(
			"ip", "1.2.3.4",
			"port", "1111",
			"ip", "1.2.3.4",
			"port", "1111",
			"ip", "1.2.3.4",
			"port", "1111",
		)
		for n := 0; n < b.N; n++ {
			ln.Printf("count: %d", b.N)
		}
	})
	b.Run("AddPrint20", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			log.Error.Add(
				"ip", "1.2.3.4",
				"port", "1111",
				"ip", "1.2.3.4",
				"port", "1111",
				"ip", "1.2.3.4",
				"port", "1111",
			).Printf("count: %d", b.N)
		}
	})
}

func ExamplePrintf() {
	log.SetOutput(os.Stdout)
	log.Service = "ex"
	log.Time = func() interface{} { return 1000 }
	log.Printf("hello, world")
	// Output: {"svc":"ex", "ts":1000, "level":"info", "msg":"hello, world"}
}

func Example() {
	log.SetOutput(os.Stdout)
	log.Service = "ex"
	log.Time = func() interface{} { return 1000 }

	log.Error.F("hello, error: %v", io.EOF)
	// Output: {"svc":"ex", "ts":1000, "level":"error", "msg":"hello, error: EOF"}
}

func Example_second() {
	log.SetOutput(os.Stdout)
	log.Service = "ex"
	log.Time = func() interface{} { return 1000 }

	log.Error.Add("severity", "high").Printf("hello, error: %v", io.EOF)
	// Output: {"svc":"ex", "ts":1000, "level":"error", "severity":"high", "msg":"hello, error: EOF"}
}

func Example_third() {
	log.SetOutput(os.Stdout)
	log.Service = "ex"
	log.Time = func() interface{} { return "2121.12.04" }

	log.Error.Add(
		"env", "prod",
		"burning", true,
		"pi", 3.14,
	).Printf("error: %v", io.EOF)
	// Output: {"svc":"ex", "ts":"2121.12.04", "level":"error", "env":"prod", "burning":true, "pi":3.14, "msg":"error: EOF"}
}
