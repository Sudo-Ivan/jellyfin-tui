package tui

import "testing"

func TestRectContains(t *testing.T) {
	r := Rect{X: 2, Y: 3, W: 4, H: 5}
	if !r.Contains(2, 3) || !r.Contains(5, 7) {
		t.Fatal("inside")
	}
	if r.Contains(1, 3) || r.Contains(6, 3) || r.Contains(2, 8) {
		t.Fatal("outside")
	}
}

func TestRectSplitX(t *testing.T) {
	r := Rect{X: 0, Y: 0, W: 10, H: 3}
	left, right := r.SplitX(4)
	if left.W != 4 || right.X != 4 || right.W != 6 {
		t.Fatalf("left=%+v right=%+v", left, right)
	}
	left, right = r.SplitX(-3)
	if left.W != 7 || right.W != 3 {
		t.Fatalf("from right left=%+v right=%+v", left, right)
	}
}

func TestRectSplitYEmpty(t *testing.T) {
	r := Rect{W: 8, H: 2}
	top, bot := r.SplitY(0)
	if !top.Empty() || bot.H != 2 {
		t.Fatalf("top=%+v bot=%+v", top, bot)
	}
}
