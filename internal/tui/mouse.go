package tui

import "strconv"

func decodeMouseSGR(b []byte) (Event, []byte, bool) {
	if len(b) < sgrMinLen || b[0] != asciiESC || b[1] != '[' || b[2] != '<' {
		return Event{}, b, false
	}
	i := sgrPrefixLen
	for i < len(b) && b[i] != 'M' && b[i] != 'm' {
		i++
	}
	if i >= len(b) {
		return Event{}, b, false
	}
	press := b[i] == 'M'
	rest := b[i+1:]
	p, ok := parseMouseParts(string(b[sgrPrefixLen:i]))
	if !ok {
		return Event{Kind: KindKey, Key: KeyEsc}, rest, true
	}
	return mouseFromSGR(p.btn, p.x-1, p.y-1, press), rest, true
}

type sgrParts struct {
	btn, x, y int
}

func parseMouseParts(s string) (sgrParts, bool) {
	a, rest, ok1 := cutInt(s)
	b, rest, ok2 := cutInt(rest)
	c, _, ok3 := cutLast(rest)
	if !ok1 || !ok2 || !ok3 {
		return sgrParts{}, false
	}
	return sgrParts{btn: a, x: b, y: c}, true
}

func cutInt(s string) (int, string, bool) {
	i := 0
	for i < len(s) && s[i] != ';' {
		i++
	}
	if i == 0 || i >= len(s) {
		return 0, s, false
	}
	n, err := strconv.Atoi(s[:i])
	if err != nil {
		return 0, s, false
	}
	return n, s[i+1:], true
}

func cutLast(s string) (int, string, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, s, false
	}
	return n, "", true
}

func mouseFromSGR(btn, x, y int, press bool) Event {
	ev := Event{Kind: KindMouse, X: x, Y: y, Press: press}
	if btn&mouseMotion != 0 {
		ev.Btn = MouseMove
		return ev
	}
	if btn&mouseWheel != 0 {
		ev.Press = true
		if btn&1 == 0 {
			ev.Btn = MouseWheelUp
		} else {
			ev.Btn = MouseWheelDown
		}
		return ev
	}
	switch btn & 3 {
	case 1:
		ev.Btn = MouseMiddle
	case 2:
		ev.Btn = MouseRight
	default:
		ev.Btn = MouseLeft
	}
	return ev
}

const (
	mouseMotion  = 32
	mouseWheel   = 64
	sgrMinLen    = 9
	sgrPrefixLen = 3
)
