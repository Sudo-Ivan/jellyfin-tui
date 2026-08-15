package tui

import "time"

const (
	enterAlt = "\x1b[?1049h\x1b[?25l\x1b[?7l\x1b[?1000h\x1b[?1006h\x1b[2J\x1b[H"
	leaveAlt = "\x1b[?1006l\x1b[?1000l\x1b[?7h\x1b[?25h\x1b[?1049l"

	eventQueue      = 64
	keyReadBuf      = 64
	escWait         = 12 * time.Millisecond
	presentCap      = 4096
	fullRedrawRatio = 0.62
	mixMax          = 256
	defaultCols     = 80
	defaultRows     = 24
	winRecCount     = 8
	inputEventSize  = 16
	escMaxPending   = 8

	sgrReset = "\x1b[0m"
	cupHome  = "\x1b[H\x1b[0m"

	asciiDEL = 127
	asciiESC = 0x1b
	asciiMin = 32
	asciiETX = 0x03
	asciiEOT = 0x04
	asciiFF  = 0x0c
	asciiNAK = 0x15
	asciiETB = 0x17
	asciiTAB = 0x09
	asciiCR  = 0x0d
	asciiLF  = 0x0a
	asciiBS  = 0x08
	asciiSP  = 0x20

	byteMask = 0xFF
	byteMax  = 255

	sgrFgTrue = "38;2;"
	sgrBgTrue = "48;2;"
)
