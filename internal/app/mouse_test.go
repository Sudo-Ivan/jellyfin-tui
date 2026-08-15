package app

import (
	"testing"
	"time"

	"jellyfin-tui/internal/tui"
)

const (
	durMMSS   = 90 * time.Second
	missCoord = 50
	pickIdx   = 3
	wantZero  = "0:00"
	wantMMSS  = "1:30"
	wantHMS   = "1:01:01"
)

func TestPickHitPrefersLaterZone(t *testing.T) {
	zones := []hitZone{
		{r: tui.Rect{X: 0, Y: 0, W: 10, H: 5}, kind: hitRail, idx: 0},
		{r: tui.Rect{X: 2, Y: 1, W: 4, H: 1}, kind: hitHome, idx: pickIdx},
	}
	z, ok := pickHit(zones, pickIdx, 1)
	if !ok || z.kind != hitHome || z.idx != pickIdx {
		t.Fatalf("%+v ok=%v", z, ok)
	}
	z, ok = pickHit(zones, 0, 0)
	if !ok || z.kind != hitRail {
		t.Fatalf("rail %+v", z)
	}
	if _, ok = pickHit(zones, missCoord, missCoord); ok {
		t.Fatal("miss")
	}
}

func TestFmtDur(t *testing.T) {
	if fmtDur(-time.Second) != wantZero {
		t.Fatal("neg")
	}
	if got := fmtDur(durMMSS); got != wantMMSS {
		t.Fatalf("mmss %s", got)
	}
	hourish := time.Duration(secondsPerHour+secondsPerMin+1) * time.Second
	if got := fmtDur(hourish); got != wantHMS {
		t.Fatalf("hms %s", got)
	}
}
