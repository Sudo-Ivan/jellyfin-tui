package jellyfin

import "testing"

func TestNextEpisodeAdjacentOrder(t *testing.T) {
	list := []Item{
		{ID: idA, ParentIndexNumber: 1, IndexNumber: 1},
		{ID: idB, ParentIndexNumber: 1, IndexNumber: 2},
		{ID: idC, ParentIndexNumber: 2, IndexNumber: 1},
	}
	next, ok := NextEpisode(list, idA)
	if !ok || next.ID != idB {
		t.Fatalf("got %+v ok=%v", next, ok)
	}
	next, ok = NextEpisode(list, idB)
	if !ok || next.ID != idC {
		t.Fatalf("season wrap got %+v ok=%v", next, ok)
	}
	_, ok = NextEpisode(list, idC)
	if ok {
		t.Fatal("last episode should have no next")
	}
}

func TestNextEpisodeUnsortedUsesSeasonIndex(t *testing.T) {
	list := []Item{
		{ID: idC, ParentIndexNumber: 2, IndexNumber: 1},
		{ID: idA, ParentIndexNumber: 1, IndexNumber: 1},
		{ID: idB, ParentIndexNumber: 1, IndexNumber: 2},
	}
	next, ok := NextEpisode(list, idB)
	if !ok || next.ID != idC {
		t.Fatalf("got %+v ok=%v", next, ok)
	}
}

func TestNextEpisodeUnknown(t *testing.T) {
	_, ok := NextEpisode([]Item{{ID: idA}}, "missing")
	if ok {
		t.Fatal("missing id")
	}
}
