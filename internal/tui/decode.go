package tui

import "unicode/utf8"

func decodeKey(b []byte) (Event, []byte, bool) {
	if len(b) == 0 {
		return Event{}, b, false
	}
	if ev, rest, ok := decodeCtrl(b); ok {
		return ev, rest, true
	}
	if b[0] == asciiESC {
		return decodeEsc(b)
	}
	if b[0] < asciiMin {
		return Event{Kind: KindKey, Key: KeyRune, Rune: rune(b[0]), Mod: ModCtrl}, b[1:], true
	}
	r, size := utf8.DecodeRune(b)
	if r == utf8.RuneError && size == 1 {
		return Event{}, b, false
	}
	return Event{Kind: KindKey, Key: KeyRune, Rune: r}, b[size:], true
}

func decodeCtrl(b []byte) (Event, []byte, bool) {
	switch b[0] {
	case asciiETX:
		return Event{Kind: KindKey, Key: KeyCtrlC}, b[1:], true
	case asciiEOT:
		return Event{Kind: KindKey, Key: KeyCtrlD}, b[1:], true
	case asciiFF:
		return Event{Kind: KindKey, Key: KeyCtrlL}, b[1:], true
	case asciiNAK:
		return Event{Kind: KindKey, Key: KeyCtrlU}, b[1:], true
	case asciiETB:
		return Event{Kind: KindKey, Key: KeyCtrlW}, b[1:], true
	case asciiTAB:
		return Event{Kind: KindKey, Key: KeyTab}, b[1:], true
	case asciiCR, asciiLF:
		return Event{Kind: KindKey, Key: KeyEnter}, b[1:], true
	case asciiDEL, asciiBS:
		return Event{Kind: KindKey, Key: KeyBackspace}, b[1:], true
	case asciiSP:
		return Event{Kind: KindKey, Key: KeySpace, Rune: ' '}, b[1:], true
	default:
		return Event{}, b, false
	}
}

func decodeEsc(b []byte) (Event, []byte, bool) {
	if len(b) == 1 {
		return Event{}, b, false
	}
	if b[1] == '[' {
		return decodeCSI(b)
	}
	if b[1] == 'O' && len(b) >= 3 {
		switch b[2] {
		case 'P':
			return Event{Kind: KindKey, Key: KeyF1}, b[3:], true
		case 'H':
			return Event{Kind: KindKey, Key: KeyHome}, b[3:], true
		case 'F':
			return Event{Kind: KindKey, Key: KeyEnd}, b[3:], true
		}
	}
	if b[1] == asciiESC {
		return Event{Kind: KindKey, Key: KeyEsc}, b[2:], true
	}
	return Event{Kind: KindKey, Key: KeyRune, Rune: rune(b[1]), Mod: ModAlt}, b[2:], true
}

func decodeCSI(b []byte) (Event, []byte, bool) {
	if len(b) < 3 {
		return Event{}, b, false
	}
	if b[2] == '<' {
		return decodeMouseSGR(b)
	}
	if ev, ok := csiArrow(b[2]); ok {
		return Event{Kind: KindKey, Key: ev}, b[3:], true
	}
	i := 2
	for i < len(b) && ((b[i] >= '0' && b[i] <= '9') || b[i] == ';') {
		i++
	}
	if i >= len(b) {
		return Event{}, b, false
	}
	final := b[i]
	body := b[2:i]
	rest := b[i+1:]
	if final == '~' {
		return csiTilde(string(body), rest)
	}
	if k, ok := csiArrow(final); ok {
		return Event{Kind: KindKey, Key: k, Mod: csiMod(body)}, rest, true
	}
	return Event{Kind: KindKey, Key: KeyEsc}, rest, true
}

func csiArrow(b byte) (Key, bool) {
	switch b {
	case 'A':
		return KeyUp, true
	case 'B':
		return KeyDown, true
	case 'C':
		return KeyRight, true
	case 'D':
		return KeyLeft, true
	case 'H':
		return KeyHome, true
	case 'F':
		return KeyEnd, true
	case 'Z':
		return KeyBacktab, true
	default:
		return KeyNone, false
	}
}

func csiTilde(body string, rest []byte) (Event, []byte, bool) {
	switch body {
	case "1", "7":
		return Event{Kind: KindKey, Key: KeyHome}, rest, true
	case "3":
		return Event{Kind: KindKey, Key: KeyDelete}, rest, true
	case "2":
		return Event{Kind: KindKey, Key: KeyInsert}, rest, true
	case "4", "8":
		return Event{Kind: KindKey, Key: KeyEnd}, rest, true
	case "5":
		return Event{Kind: KindKey, Key: KeyPgUp}, rest, true
	case "6":
		return Event{Kind: KindKey, Key: KeyPgDn}, rest, true
	default:
		return Event{Kind: KindKey, Key: KeyEsc}, rest, true
	}
}

func csiMod(_ []byte) Mod {
	return 0
}
