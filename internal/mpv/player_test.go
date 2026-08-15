package mpv

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

const (
	errDial     = "connect: no such file or directory"
	errInterp   = "Interpreter not found!"
	wantIPCWrap = "mpv ipc:"
)

func TestAbortPlayerIncludesStderr(t *testing.T) {
	var buf bytes.Buffer
	if _, err := buf.WriteString(errInterp + "\n"); err != nil {
		t.Fatal(err)
	}
	err := abortPlayer(&exec.Cmd{}, "", &buf, fmt.Errorf("%s", errDial))
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if !strings.Contains(msg, wantIPCWrap) || !strings.Contains(msg, errInterp) || !strings.Contains(msg, errDial) {
		t.Fatalf("err %q", msg)
	}
}

func TestAbortPlayerWithoutStderr(t *testing.T) {
	err := abortPlayer(&exec.Cmd{}, "", &bytes.Buffer{}, fmt.Errorf("%s", errDial))
	if err == nil || !strings.Contains(err.Error(), errDial) {
		t.Fatalf("err %v", err)
	}
	if strings.Contains(err.Error(), "(") {
		t.Fatalf("unexpected stderr wrap %v", err)
	}
}
