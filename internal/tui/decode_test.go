package tui

import "testing"

func TestDecodeEnterAndCtrlC(t *testing.T) {
	ev, rest, ok := decodeKey([]byte{asciiCR, 'x'})
	if !ok || ev.Key != KeyEnter || len(rest) != 1 {
		t.Fatalf("%+v rest=%q ok=%v", ev, rest, ok)
	}
	ev, _, ok = decodeKey([]byte{asciiETX})
	if !ok || ev.Key != KeyCtrlC {
		t.Fatal("ctrl-c")
	}
}

func TestDecodeCSIArrows(t *testing.T) {
	ev, _, ok := decodeKey([]byte{asciiESC, '[', 'A'})
	if !ok || ev.Key != KeyUp {
		t.Fatalf("%+v", ev)
	}
}
