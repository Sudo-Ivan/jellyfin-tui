package tui

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	enableProcessedInput = 0x0001
	enableLineInput      = 0x0002
	enableEchoInput      = 0x0004
	enableWindowInput    = 0x0008
	enableMouseInput     = 0x0010
	enableQuickEdit      = 0x0040
	enableVTInput        = 0x0200
	enableProcessedOut   = 0x0001
	enableWrapAtEOL      = 0x0002
	enableVTProc         = 0x0004
	enableVTDisableNL    = 0x0008

	stdInputHandle  = ^uintptr(10 - 1)
	stdOutputHandle = ^uintptr(11 - 1)
)

var (
	k32                        = syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle           = k32.NewProc("GetStdHandle")
	procGetConsoleMode         = k32.NewProc("GetConsoleMode")
	procSetConsoleMode         = k32.NewProc("SetConsoleMode")
	procGetConsoleScreenBuffer = k32.NewProc("GetConsoleScreenBufferInfo")
	procReadConsoleInput       = k32.NewProc("ReadConsoleInputW")
)

type coord struct{ X, Y int16 }

type smallRect struct{ Left, Top, Right, Bottom int16 }

type consoleScreenBufferInfo struct {
	Size       coord
	Cursor     coord
	Attrs      uint16
	Window     smallRect
	MaxWinSize coord
}

type winSaved struct {
	inMode  uint32
	outMode uint32
}

func Open() (*Term, error) {
	in := os.Stdin
	out := os.Stdout
	hin, err := getStdHandle(stdInputHandle)
	if err != nil {
		return nil, err
	}
	hout, err := getStdHandle(stdOutputHandle)
	if err != nil {
		return nil, err
	}
	inMode, err := getConsoleMode(hin)
	if err != nil {
		return nil, err
	}
	outMode, err := getConsoleMode(hout)
	if err != nil {
		return nil, err
	}
	rawIn := (inMode | enableWindowInput | enableVTInput | enableMouseInput) &^
		(enableLineInput | enableEchoInput | enableProcessedInput | enableQuickEdit)
	rawOut := outMode | enableProcessedOut | enableVTProc | enableVTDisableNL
	if err := setConsoleMode(hin, rawIn); err != nil {
		return nil, err
	}
	if err := setConsoleMode(hout, rawOut); err != nil {
		_ = setConsoleMode(hin, inMode)
		return nil, err
	}
	t := &Term{
		in:    in,
		out:   out,
		theme: Oxide(),
		saved: winSaved{inMode: inMode, outMode: outMode},
	}
	w, h, err := t.readSize()
	if err != nil {
		_ = t.Restore()
		return nil, err
	}
	t.w, t.h = w, h
	t.buf = NewBuffer(w, h)
	if err := t.write(enterAlt); err != nil {
		_ = t.Restore()
		return nil, err
	}
	return t, nil
}

func (t *Term) Restore() error {
	_ = t.write(leaveAlt)
	hin, err := getStdHandle(stdInputHandle)
	if err != nil {
		return err
	}
	hout, err := getStdHandle(stdOutputHandle)
	if err != nil {
		return err
	}
	if s, ok := t.saved.(winSaved); ok {
		_ = setConsoleMode(hin, s.inMode)
		_ = setConsoleMode(hout, s.outMode)
	}
	return nil
}

func (t *Term) readSize() (cols int, rows int, err error) {
	hout, err := getStdHandle(stdOutputHandle)
	if err != nil {
		return 0, 0, err
	}
	var info consoleScreenBufferInfo
	r1, _, e := procGetConsoleScreenBuffer.Call(uintptr(hout), uintptr(unsafe.Pointer(&info))) // #nosec G103
	if r1 == 0 {
		if e != syscall.Errno(0) {
			return 0, 0, e
		}
		return defaultCols, defaultRows, nil
	}
	w := int(info.Window.Right - info.Window.Left + 1)
	h := int(info.Window.Bottom - info.Window.Top + 1)
	if w < 1 {
		w = defaultCols
	}
	if h < 1 {
		h = defaultRows
	}
	return w, h, nil
}

func getStdHandle(n uintptr) (syscall.Handle, error) {
	r1, _, e := procGetStdHandle.Call(n)
	if r1 == 0 || r1 == uintptr(^uintptr(0)) {
		if e != syscall.Errno(0) {
			return 0, e
		}
		return 0, syscall.EINVAL
	}
	return syscall.Handle(r1), nil
}

func getConsoleMode(h syscall.Handle) (uint32, error) {
	var mode uint32
	r1, _, e := procGetConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode))) // #nosec G103
	if r1 == 0 {
		return 0, e
	}
	return mode, nil
}

func setConsoleMode(h syscall.Handle, mode uint32) error {
	r1, _, e := procSetConsoleMode.Call(uintptr(h), uintptr(mode))
	if r1 == 0 {
		return e
	}
	return nil
}
