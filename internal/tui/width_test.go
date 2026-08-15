package tui

import (
	"strings"
	"testing"
)

const (
	wideHan = "中文标题"
	wideMix = "ab中文"
)

func TestEllipsizeWideRunes(t *testing.T) {
	if ellipsize(wideHan, 8) != wideHan {
		t.Fatalf("fits %q", ellipsize(wideHan, 8))
	}
	got := ellipsize(wideHan, 3)
	if visWidth(got) > 3 {
		t.Fatalf("width %q", got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("ellipsis %q", got)
	}
}

func TestClipVisAtBoundary(t *testing.T) {
	if clipVis(wideMix, 3) != "ab" {
		t.Fatalf("clip %q", clipVis(wideMix, 3))
	}
	if clipVis("hello", 0) != "" {
		t.Fatal("zero")
	}
}

func TestCombiningMarkWidth(t *testing.T) {
	const eAcute = "e\u0301"
	if RuneWidth('\u0301') != 0 {
		t.Fatal("combining width")
	}
	if visWidth(eAcute) != 1 {
		t.Fatalf("vis %d", visWidth(eAcute))
	}
}

func TestEllipsizeSingleColumn(t *testing.T) {
	if ellipsize("hello", 1) != "h" {
		t.Fatalf("%q", ellipsize("hello", 1))
	}
}
