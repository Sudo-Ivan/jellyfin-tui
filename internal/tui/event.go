package tui

import "time"

// Kind classifies an input event.
type Kind uint8

const (
	KindKey Kind = iota
	KindResize
	KindTick
	KindMouse
)

// Key is a non-character key or a character wrapped as KeyRune.
type Key uint16

const (
	KeyNone Key = iota
	KeyRune
	KeyEnter
	KeyEsc
	KeyBackspace
	KeyTab
	KeyBacktab
	KeySpace
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDn
	KeyDelete
	KeyInsert
	KeyF1
	KeyF2
	KeyCtrlC
	KeyCtrlD
	KeyCtrlL
	KeyCtrlU
	KeyCtrlW
)

// Mod is a bitset of modifiers.
type Mod uint8

const (
	ModCtrl Mod = 1 << iota
	ModAlt
	ModShift
)

// MouseBtn is a mouse button or wheel direction.
type MouseBtn uint8

const (
	MouseNone MouseBtn = iota
	MouseLeft
	MouseMiddle
	MouseRight
	MouseWheelUp
	MouseWheelDown
	MouseMove
)

// Event is a single input or window change.
type Event struct {
	Kind  Kind
	Key   Key
	Rune  rune
	Mod   Mod
	W, H  int
	X, Y  int
	Btn   MouseBtn
	Press bool
	At    time.Time
}

func (e Event) RuneKey(r rune) bool {
	return e.Kind == KindKey && e.Key == KeyRune && e.Rune == r && e.Mod == 0
}

func (e Event) Is(k Key) bool {
	return e.Kind == KindKey && e.Key == k
}

func (e Event) Click() bool {
	return e.Kind == KindMouse && e.Btn == MouseLeft && e.Press
}

func (e Event) Wheel() int {
	if e.Kind != KindMouse {
		return 0
	}
	switch e.Btn {
	case MouseWheelUp:
		return -1
	case MouseWheelDown:
		return 1
	default:
		return 0
	}
}
