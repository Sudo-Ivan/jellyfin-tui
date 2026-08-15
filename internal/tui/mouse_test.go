package tui

import "testing"

const (
	sgrCol = 5
	sgrRow = 10
)

func TestDecodeMouseSGRPress(t *testing.T) {
	b := []byte{asciiESC, '[', '<', '0', ';', '5', ';', '1', '0', 'M'}
	ev, rest, ok := decodeKey(b)
	if !ok || ev.Kind != KindMouse || ev.Btn != MouseLeft || !ev.Press {
		t.Fatalf("%+v ok=%v", ev, ok)
	}
	if ev.X != sgrCol-1 || ev.Y != sgrRow-1 || len(rest) != 0 {
		t.Fatalf("pos %d,%d rest=%q", ev.X, ev.Y, rest)
	}
	if !ev.Click() {
		t.Fatal("click")
	}
}

func TestDecodeMouseSGRRelease(t *testing.T) {
	b := []byte{asciiESC, '[', '<', '0', ';', '1', ';', '1', 'm'}
	ev, _, ok := decodeKey(b)
	if !ok || ev.Press || ev.Click() {
		t.Fatalf("%+v", ev)
	}
}

func TestDecodeMouseSGRWheel(t *testing.T) {
	up := []byte{asciiESC, '[', '<', '6', '4', ';', '1', ';', '1', 'M'}
	ev, _, ok := decodeKey(up)
	if !ok || ev.Wheel() != -1 {
		t.Fatalf("up %+v", ev)
	}
	down := []byte{asciiESC, '[', '<', '6', '5', ';', '1', ';', '1', 'M'}
	ev, _, ok = decodeKey(down)
	if !ok || ev.Wheel() != 1 {
		t.Fatalf("down %+v", ev)
	}
}

func TestDecodeMouseSGRIncomplete(t *testing.T) {
	_, _, ok := decodeKey([]byte{asciiESC, '[', '<', '0', ';'})
	if ok {
		t.Fatal("want wait")
	}
}
