// Package tui is a retained-cell terminal UI with dirty-cell present.
package tui

import (
	"io"
	"os"
)

// Term owns the host tty in raw mode plus the retained raster.
type Term struct {
	in    *os.File
	out   *os.File
	buf   *Buffer
	theme Theme
	w, h  int
	saved any
}

func (t *Term) Size() (w, h int) { return t.w, t.h }
func (t *Term) Out() io.Writer   { return t.out }
func (t *Term) Buffer() *Buffer  { return t.buf }

func (t *Term) Painter() *Painter {
	return NewPainter(t.buf, t.theme)
}

func (t *Term) SetTheme(th Theme) { t.theme = th }

func (t *Term) write(s string) error {
	_, err := t.out.Write([]byte(s))
	return err
}

func (t *Term) Present() error {
	return t.buf.Present(t.out)
}

func (t *Term) Sync() error {
	t.buf.Invalidate()
	return t.Present()
}

func (t *Term) Resize() error {
	w, h, err := t.readSize()
	if err != nil {
		return err
	}
	if w == t.w && h == t.h {
		return nil
	}
	t.w, t.h = w, h
	t.buf.Resize(w, h)
	return nil
}
