package mpv

import (
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testPosSec    = 12.5
	testVolLevel  = 80
	testPath      = "/media/movie.mkv"
	testDurMillis = 2500
)

func TestApplyPropertyChange(t *testing.T) {
	p := &Player{}
	p.apply(map[string]any{
		"event": evProp,
		"name":  propTime,
		"data":  testPosSec,
	})
	st := p.Status()
	if st.Pos != time.Duration(testPosSec*float64(time.Second)) {
		t.Fatalf("pos %v", st.Pos)
	}
	p.apply(map[string]any{
		"event": evProp,
		"name":  propVolume,
		"data":  json.Number("80"),
	})
	if p.Status().Volume != testVolLevel {
		t.Fatalf("vol %d", p.Status().Volume)
	}
	p.apply(map[string]any{
		"event": evProp,
		"name":  propPath,
		"data":  testPath,
	})
	if p.Status().Path != testPath {
		t.Fatalf("path %q", p.Status().Path)
	}
}

func TestApplyEndFileEOFCallsCallback(t *testing.T) {
	var called atomic.Bool
	p := &Player{onEOF: func() { called.Store(true) }}
	p.apply(map[string]any{"event": evEndFile, "reason": reasonEOF})
	if !p.Status().EOF {
		t.Fatal("eof flag")
	}
	if !called.Load() {
		t.Fatal("callback")
	}
}

func TestApplyEndFileStopNoCallback(t *testing.T) {
	var called atomic.Bool
	p := &Player{onEOF: func() { called.Store(true) }}
	p.apply(map[string]any{"event": evEndFile, "reason": reasonStop})
	if p.Status().EOF {
		t.Fatal("stop is not eof")
	}
	if called.Load() {
		t.Fatal("callback on stop")
	}
}

func TestApplyIgnoresUnknownEvent(t *testing.T) {
	p := &Player{vol: defaultVolume}
	p.apply(map[string]any{"event": "file-loaded"})
	if p.Status().Volume != defaultVolume {
		t.Fatal("unchanged")
	}
}

func TestAsDurationTypes(t *testing.T) {
	if asDuration(nil) != 0 {
		t.Fatal("nil")
	}
	if asDuration(json.Number("2.5")) != testDurMillis*time.Millisecond {
		t.Fatal("json.Number")
	}
	if asInt(json.Number("7")) != 7 {
		t.Fatal("int")
	}
}
