package log

import (
	"io"
	"testing"
)

func TestLog(t *testing.T) {
	Service = "testlog"
	Time = func() interface{} { return 12345 }
	have := line{}.Sprintf("test log message: %s", "package scoped variables arent hard to test")
	want := `{"svc":"testlog", "time":12345, "level":"info", "msg":"test log message: package scoped variables arent hard to test"}`
	if have != want {
		t.Fatalf("bad log:\n\t\thave: %s\n\t\twant: %s", have, want)
	}
}

func TestFatal(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("didnt panic")
		}
	}()
	Fatal("panic: %v", io.EOF)
}
