package app

import "jellyfin-tui/internal/tui"

type field struct {
	runes []rune
	cur   int
}

func (f *field) set(s string) {
	f.runes = []rune(s)
	f.cur = len(f.runes)
}

func (f *field) string() string { return string(f.runes) }

func (f *field) handle(e tui.Event) {
	switch {
	case e.Is(tui.KeyLeft) && f.cur > 0:
		f.cur--
	case e.Is(tui.KeyRight) && f.cur < len(f.runes):
		f.cur++
	case e.Is(tui.KeyHome):
		f.cur = 0
	case e.Is(tui.KeyEnd):
		f.cur = len(f.runes)
	case e.Is(tui.KeyBackspace):
		f.backspace()
	case e.Is(tui.KeyDelete):
		f.delete()
	case e.Is(tui.KeyCtrlU):
		f.runes = nil
		f.cur = 0
	case e.Is(tui.KeyCtrlW):
		f.deleteWord()
	case e.Key == tui.KeyRune && e.Mod == 0 && e.Rune >= printableMin:
		f.insert(e.Rune)
	case e.Is(tui.KeySpace):
		f.insert(' ')
	}
}

func (f *field) backspace() {
	if f.cur == 0 {
		return
	}
	f.runes = append(f.runes[:f.cur-1], f.runes[f.cur:]...)
	f.cur--
}

func (f *field) delete() {
	if f.cur < len(f.runes) {
		f.runes = append(f.runes[:f.cur], f.runes[f.cur+1:]...)
	}
}

func (f *field) deleteWord() {
	i := f.cur
	for i > 0 && f.runes[i-1] == ' ' {
		i--
	}
	for i > 0 && f.runes[i-1] != ' ' {
		i--
	}
	f.runes = append(f.runes[:i], f.runes[f.cur:]...)
	f.cur = i
}

func (f *field) insert(r rune) {
	f.runes = append(f.runes, 0)
	copy(f.runes[f.cur+1:], f.runes[f.cur:])
	f.runes[f.cur] = r
	f.cur++
}

type fieldPaint struct {
	label  string
	value  string
	cur    int
	secret bool
	focus  bool
}

func paintField(p *tui.Painter, r tui.Rect, f fieldPaint, th tui.Theme) {
	p.Text(r.X, r.Y, f.label, th.Muted().WithBg(th.Basin))
	box := tui.Rect{X: r.X + fieldLabelW, Y: r.Y, W: r.W - fieldLabelW - 2, H: 1}
	st := th.Panel()
	if f.focus {
		st = th.Cursor()
	}
	p.Fill(box, st)
	shown := f.value
	if f.secret {
		b := make([]rune, len([]rune(f.value)))
		for i := range b {
			b[i] = secretCh
		}
		shown = string(b)
	}
	p.TextClip(tui.Rect{X: box.X + 1, Y: box.Y, W: box.W - 2, H: 1}, shown, st)
	if !f.focus {
		return
	}
	cx := box.X + 1 + f.cur
	if f.secret {
		n := len([]rune(f.value))
		cx = box.X + 1 + n
		if f.cur < n {
			cx = box.X + 1 + f.cur
		}
	}
	if cx >= box.X+1 && cx < box.MaxX()-1 {
		p.Put(cx, box.Y, caretCh, st.WithFg(th.Ember))
	}
}
