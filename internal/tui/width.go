package tui

var combining = [][2]rune{
	{0x0300, 0x036F},
	{0x20D0, 0x20FF},
	{0xFE20, 0xFE2F},
}

var wide = [][2]rune{
	{0x1100, 0x115F},
	{0x2E80, 0xA4CF},
	{0xAC00, 0xD7A3},
	{0xF900, 0xFAFF},
	{0xFE10, 0xFE19},
	{0xFE30, 0xFE6F},
	{0xFF00, 0xFF60},
	{0xFFE0, 0xFFE6},
	{0x1F300, 0x1FAFF},
}

const (
	angleL     = 0x2329
	angleR     = 0x232A
	ideographN = 0x303F
)

func RuneWidth(r rune) int {
	switch {
	case r == 0 || r == asciiDEL:
		return 0
	case r < asciiMin:
		return 0
	case r < asciiDEL:
		return 1
	case inRanges(r, combining):
		return 0
	case wideRune(r):
		return 2
	default:
		return 1
	}
}

func wideRune(r rune) bool {
	if r == angleL || r == angleR {
		return true
	}
	if r == ideographN {
		return false
	}
	return inRanges(r, wide)
}

func inRanges(r rune, ranges [][2]rune) bool {
	for _, rg := range ranges {
		if r >= rg[0] && r <= rg[1] {
			return true
		}
	}
	return false
}

func visWidth(s string) int {
	n := 0
	for _, r := range s {
		n += RuneWidth(r)
	}
	return n
}

func clipVis(s string, limit int) string {
	if limit <= 0 {
		return ""
	}
	n := 0
	for i, r := range s {
		w := RuneWidth(r)
		if n+w > limit {
			return s[:i]
		}
		n += w
	}
	return s
}

func ellipsize(s string, limit int) string {
	if limit <= 0 {
		return ""
	}
	if visWidth(s) <= limit {
		return s
	}
	if limit == 1 {
		return clipVis(s, 1)
	}
	return clipVis(s, limit-1) + "…"
}
