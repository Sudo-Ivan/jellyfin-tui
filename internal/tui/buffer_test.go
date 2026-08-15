package tui

import (
	"bytes"
	"strings"
	"testing"
)

const testWhite byte = 0xFF

var (
	testFg = RGB(testWhite, testWhite, testWhite)
	testBg = RGB(0, 0, 0)
)

func TestDirtyPresentWritesOnlyChangedCells(t *testing.T) {
	b := NewBuffer(8, 2)
	st := Style{Fg: testFg, Bg: testBg}
	b.Clear(st)
	var first bytes.Buffer
	if err := b.Present(&first); err != nil {
		t.Fatal(err)
	}
	if first.Len() == 0 {
		t.Fatal("first present wrote nothing")
	}
	b.Put(3, 0, Cell{Ch: 'A', Style: st})
	var second bytes.Buffer
	if err := b.Present(&second); err != nil {
		t.Fatal(err)
	}
	out := second.String()
	if !strings.Contains(out, "A") {
		t.Fatalf("expected dirty cell in output, got %q", out)
	}
	if strings.Count(out, "A") != 1 {
		t.Fatalf("expected a single dirty glyph, got %q", out)
	}
}

func TestResizeForcesFullRedraw(t *testing.T) {
	b := NewBuffer(4, 1)
	b.Clear(Style{})
	var buf bytes.Buffer
	if err := b.Present(&buf); err != nil {
		t.Fatal(err)
	}
	b.Resize(6, 2)
	if !b.force {
		t.Fatal("resize should invalidate the back buffer")
	}
}

func TestRuneWidthCJK(t *testing.T) {
	if RuneWidth('A') != 1 {
		t.Fatal("ascii")
	}
	if RuneWidth('字') != 2 {
		t.Fatal("cjk")
	}
	if visWidth("字A") != 3 {
		t.Fatal("mixed")
	}
}
