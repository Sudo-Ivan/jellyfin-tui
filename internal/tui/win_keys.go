//go:build windows

package tui

const (
	keyEventType    = 0x0001
	mouseEventType  = 0x0002
	windowSizeEvent = 0x0004

	mouseMoved   = 0x0001
	mouseWheeled = 0x0004
	leftButton   = 0x0001

	vkBack   = 0x08
	vkTab    = 0x09
	vkReturn = 0x0D
	vkEsc    = 0x1B
	vkSpace  = 0x20
	vkPrior  = 0x21
	vkNext   = 0x22
	vkEnd    = 0x23
	vkHome   = 0x24
	vkLeft   = 0x25
	vkUp     = 0x26
	vkRight  = 0x27
	vkDown   = 0x28
	vkInsert = 0x2D
	vkDelete = 0x2E
	vkF1     = 0x70
	vkF2     = 0x71

	shiftPressed = 0x0010
	leftCtrl     = 0x0008
	rightCtrl    = 0x0004
	leftAlt      = 0x0002
	rightAlt     = 0x0001
)

type inputRecord struct {
	EventType uint16
	_         uint16
	Event     [inputEventSize]byte
}

type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

type mouseEventRecord struct {
	Pos        coord
	Buttons    uint32
	Control    uint32
	EventFlags uint32
}

func winKey(ke keyEventRecord) (Event, bool) {
	mod := winMod(ke.ControlKeyState)
	if k, ok := vkMap[ke.VirtualKeyCode]; ok {
		if ke.VirtualKeyCode == vkTab && mod&ModShift != 0 {
			return Event{Kind: KindKey, Key: KeyBacktab, Mod: mod}, true
		}
		if ke.VirtualKeyCode == vkSpace {
			return Event{Kind: KindKey, Key: KeySpace, Rune: ' ', Mod: mod}, true
		}
		return Event{Kind: KindKey, Key: k, Mod: mod}, true
	}
	return winRune(ke, mod)
}

func winMouse(me mouseEventRecord) (Event, bool) {
	x, y := int(me.Pos.X), int(me.Pos.Y)
	if me.EventFlags&mouseWheeled != 0 {
		btn := MouseWheelDown
		if int32(me.Buttons) > 0 {
			btn = MouseWheelUp
		}
		return Event{Kind: KindMouse, Btn: btn, X: x, Y: y, Press: true}, true
	}
	if me.EventFlags&mouseMoved != 0 {
		return Event{}, false
	}
	if me.Buttons&leftButton == 0 {
		return Event{}, false
	}
	return Event{Kind: KindMouse, Btn: MouseLeft, X: x, Y: y, Press: true}, true
}

func winMod(state uint32) Mod {
	var mod Mod
	if state&(leftCtrl|rightCtrl) != 0 {
		mod |= ModCtrl
	}
	if state&(leftAlt|rightAlt) != 0 {
		mod |= ModAlt
	}
	if state&shiftPressed != 0 {
		mod |= ModShift
	}
	return mod
}

var vkMap = map[uint16]Key{
	vkLeft:   KeyLeft,
	vkRight:  KeyRight,
	vkUp:     KeyUp,
	vkDown:   KeyDown,
	vkHome:   KeyHome,
	vkEnd:    KeyEnd,
	vkPrior:  KeyPgUp,
	vkNext:   KeyPgDn,
	vkDelete: KeyDelete,
	vkInsert: KeyInsert,
	vkReturn: KeyEnter,
	vkEsc:    KeyEsc,
	vkBack:   KeyBackspace,
	vkTab:    KeyTab,
	vkF1:     KeyF1,
	vkF2:     KeyF2,
	vkSpace:  KeySpace,
}

func winRune(ke keyEventRecord, mod Mod) (Event, bool) {
	ch := rune(ke.UnicodeChar)
	if ch == 0 {
		return Event{}, false
	}
	if mod&ModCtrl != 0 {
		if k, ok := ctrlRune[ch]; ok {
			return Event{Kind: KindKey, Key: k}, true
		}
	}
	if ch == ' ' {
		return Event{Kind: KindKey, Key: KeySpace, Rune: ' ', Mod: mod}, true
	}
	return Event{Kind: KindKey, Key: KeyRune, Rune: ch, Mod: mod}, true
}

var ctrlRune = map[rune]Key{
	3:  KeyCtrlC,
	4:  KeyCtrlD,
	12: KeyCtrlL,
	21: KeyCtrlU,
	23: KeyCtrlW,
}
