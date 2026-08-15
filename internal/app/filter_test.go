package app

import (
	"testing"

	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

const (
	yearHeat = 1995
	yearClue = 1985
)

func TestBrowseFilterShownGenreActorYear(t *testing.T) {
	heat := jellyfin.Item{
		ID: "1", Name: "Heat", ProductionYear: yearHeat,
		Genres: []string{"Crime"},
		People: []jellyfin.Person{{Name: "Pacino", Type: "Actor"}},
	}
	clue := jellyfin.Item{
		ID: "2", Name: "Clue", ProductionYear: yearClue,
		Genres: []string{"Comedy"},
	}
	heat2 := jellyfin.Item{
		ID: "3", Name: "Heat 2", ProductionYear: yearHeat,
		Genres: []string{"Crime"},
	}
	b := browseState{items: []jellyfin.Item{heat, clue, heat2}}
	b.filter.set("genre:Crime year:1995 actor:Pacino")
	b.filterShown()
	if len(b.shown) != 1 || b.items[b.shown[0]].ID != "1" {
		t.Fatalf("shown=%v", b.shown)
	}
	b.filter.set("heat")
	b.filterShown()
	if len(b.shown) != 2 {
		t.Fatalf("term shown=%v", b.shown)
	}
	b.filter.set("")
	b.filterShown()
	if len(b.shown) != 3 {
		t.Fatal("empty")
	}
}

func TestMoveSel(t *testing.T) {
	if moveSel(0, 0, tui.Event{Kind: tui.KindKey, Key: tui.KeyDown}) != 0 {
		t.Fatal("empty")
	}
	if moveSel(0, 5, tui.Event{Kind: tui.KindKey, Key: tui.KeyDown}) != 1 {
		t.Fatal("down")
	}
	if moveSel(0, 5, tui.Event{Kind: tui.KindKey, Key: tui.KeyUp}) != 0 {
		t.Fatal("up clamp")
	}
	if moveSel(2, 5, tui.Event{Kind: tui.KindKey, Key: tui.KeyEnd}) != 4 {
		t.Fatal("end")
	}
	if moveSel(4, 5, tui.Event{Kind: tui.KindKey, Key: tui.KeyHome}) != 0 {
		t.Fatal("home")
	}
}

func TestFieldEdit(t *testing.T) {
	var f field
	f.set("ab")
	f.cur = 1
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyRune, Rune: 'X'})
	if f.string() != "aXb" {
		t.Fatalf("%q", f.string())
	}
	f.handle(tui.Event{Kind: tui.KindKey, Key: tui.KeyBackspace})
	if f.string() != "ab" {
		t.Fatalf("bs %q", f.string())
	}
}
