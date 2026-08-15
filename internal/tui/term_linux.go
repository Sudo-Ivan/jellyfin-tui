package tui

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	tcgets     = 0x5401
	tcsets     = 0x5402
	tiocgwinsz = 0x5413

	ignbrk = 0x1
	brkint = 0x2
	parmrk = 0x8
	istrip = 0x20
	inlcr  = 0x40
	igncr  = 0x80
	icrnl  = 0x100
	ixon   = 0x400
	opost  = 0x1
	echo   = 0x8
	echonl = 0x40
	icanon = 0x2
	isig   = 0x1
	iexten = 0x8000
	csize  = 0x30
	parenb = 0x100
	cs8    = 0x30
)

type winSize struct {
	Row, Col, X, Y uint16
}

func Open() (*Term, error) {
	in := os.Stdin
	out := os.Stdout
	fd := int(in.Fd())
	var old syscall.Termios
	if errno := ioctlTermios(fd, tcgets, &old); errno != 0 {
		return nil, errno
	}
	raw := old
	raw.Iflag &^= ignbrk | brkint | parmrk | istrip | inlcr | igncr | icrnl | ixon
	raw.Oflag &^= opost
	raw.Lflag &^= echo | echonl | icanon | isig | iexten
	raw.Cflag &^= csize | parenb
	raw.Cflag |= cs8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if errno := ioctlTermios(fd, tcsets, &raw); errno != 0 {
		return nil, errno
	}
	t := &Term{in: in, out: out, theme: Oxide(), saved: old}
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
	if old, ok := t.saved.(syscall.Termios); ok {
		if errno := ioctlTermios(int(t.in.Fd()), tcsets, &old); errno != 0 {
			return errno
		}
	}
	return nil
}

func (t *Term) readSize() (cols int, rows int, err error) {
	var ws winSize
	if errno := ioctlWinSize(int(t.out.Fd()), &ws); errno != 0 {
		return 0, 0, errno
	}
	w, h := int(ws.Col), int(ws.Row)
	if w < 1 {
		w = defaultCols
	}
	if h < 1 {
		h = defaultRows
	}
	return w, h, nil
}

func ioctlTermios(fd int, req uintptr, tio *syscall.Termios) syscall.Errno {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(unsafe.Pointer(tio))) // #nosec G103
	return errno
}

func ioctlWinSize(fd int, ws *winSize) syscall.Errno {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tiocgwinsz, uintptr(unsafe.Pointer(ws))) // #nosec G103
	return errno
}
