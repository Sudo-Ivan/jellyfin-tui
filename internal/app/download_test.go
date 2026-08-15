package app

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

const (
	hdrName  = "X-Emby-Token"
	hdrValue = "secret-token"
)

func TestSplitHeader(t *testing.T) {
	name, val, ok := splitHeader(hdrName + ": " + hdrValue)
	if !ok || name != hdrName || val != hdrValue {
		t.Fatalf("split %q %q ok=%v", name, val, ok)
	}
	name, val, ok = splitHeader(hdrName + ":" + hdrValue)
	if !ok || val != hdrValue {
		t.Fatalf("no space %q %q", name, val)
	}
	if _, _, ok := splitHeader("nocolon"); ok {
		t.Fatal("no colon")
	}
}

func TestCopyWithProgress(t *testing.T) {
	const total = 100
	var last int
	var sink bytes.Buffer
	n, err := copyWithProgress(&sink, strings.NewReader("0123456789"), total, func(pct int) {
		last = pct
	})
	if err != nil || n != 10 {
		t.Fatalf("copy n=%d err=%v", n, err)
	}
	if last != 10 {
		t.Fatalf("pct %d", last)
	}
}

func TestCopyWithProgressEOF(t *testing.T) {
	var sink bytes.Buffer
	n, err := copyWithProgress(&sink, strings.NewReader(""), 0, nil)
	if err != nil || n != 0 {
		t.Fatalf("empty n=%d err=%v", n, err)
	}
}

func TestCopyWithProgressWriteError(t *testing.T) {
	noop := func(int) {}
	_, err := copyWithProgress(errWriter{}, strings.NewReader("x"), 1, noop)
	if err != io.ErrShortWrite {
		t.Fatalf("err %v", err)
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, io.ErrShortWrite }
