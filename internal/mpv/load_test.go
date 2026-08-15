package mpv

import (
	"testing"
	"time"
)

const (
	nameShow  = "Show"
	wantStart = "10"
)

func TestLoadOpts(t *testing.T) {
	o := loadOpts(nameShow, 10*time.Second)
	if o[optTitle] != nameShow || o[optStart] != wantStart {
		t.Fatalf("%v", o)
	}
	if len(loadOpts("", 0)) != 0 {
		t.Fatal("empty")
	}
}
