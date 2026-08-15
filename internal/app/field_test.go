package app

import (
	"testing"

	"jellyfin-tui/internal/tui"
)

func TestFieldDeleteWord(t *testing.T) {
	var f field
	f.set("hello world test")
	f.cur = len(f.runes)
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyCtrlW})
	if f.string() != "hello world " {
		t.Fatalf("ctrl-w %q", f.string())
	}
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyCtrlW})
	if f.string() != "hello " {
		t.Fatalf("ctrl-w again %q", f.string())
	}
}

func TestFieldCtrlUClears(t *testing.T) {
	var f field
	f.set("secret")
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyCtrlU})
	if f.string() != "" || f.cur != 0 {
		t.Fatalf("ctrl-u %q cur=%d", f.string(), f.cur)
	}
}

func TestFieldCursorBounds(t *testing.T) {
	var f field
	f.set("abc")
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyHome})
	if f.cur != 0 {
		t.Fatal("home")
	}
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyEnd})
	if f.cur != 3 {
		t.Fatal("end")
	}
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyLeft})
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyDelete})
	if f.string() != "ab" {
		t.Fatalf("delete %q", f.string())
	}
}
