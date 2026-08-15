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

func TestDecodeCSITildeKeys(t *testing.T) {
	ev, _, ok := decodeKey([]byte{asciiESC, '[', '3', '~'})
	if !ok || ev.Key != KeyDelete {
		t.Fatalf("%+v", ev)
	}
	ev, _, ok = decodeKey([]byte{asciiESC, '[', '5', '~'})
	if !ok || ev.Key != KeyPgUp {
		t.Fatal("pgup")
	}
}

func TestDecodeAltRune(t *testing.T) {
	ev, rest, ok := decodeKey([]byte{asciiESC, 'a', 'b'})
	if !ok || ev.Rune != 'a' || ev.Mod != ModAlt || len(rest) != 1 {
		t.Fatalf("%+v rest=%q", ev, rest)
	}
}

func TestDecodeIncompleteCSI(t *testing.T) {
	_, rest, ok := decodeKey([]byte{asciiESC, '['})
	if ok || len(rest) != 2 {
		t.Fatalf("incomplete ok=%v rest=%q", ok, rest)
	}
}
