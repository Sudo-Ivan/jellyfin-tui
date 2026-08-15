package jellyfin

import "testing"

func TestItemLabel(t *testing.T) {
	ep := Item{Type: TypeEpisode, SeriesName: "Wire", Name: "Pilot", ParentIndexNumber: 1, IndexNumber: 1}
	if ep.Label() != "Wire S1E1  Pilot" {
		t.Fatalf("%q", ep.Label())
	}
	mv := Item{Type: TypeMovie, Name: "Heat", ProductionYear: yr1995}
	if mv.Label() != "Heat  (1995)" {
		t.Fatalf("%q", mv.Label())
	}
	mv0 := Item{Type: TypeMovie, Name: "Heat"}
	if mv0.Label() != "Heat" {
		t.Fatal(mv0.Label())
	}
	sn := Item{Type: TypeSeason, IndexNumber: 2, Name: "Season 2"}
	if sn.Label() != "Season 2" {
		t.Fatal(sn.Label())
	}
}

func TestDurationAndResume(t *testing.T) {
	it := Item{RunTimeTicks: ticksPerSec * durTicks, UserData: UserData{PlaybackPositionTicks: ticksPerSec * durResume}}
	if it.Duration() != durTicks {
		t.Fatalf("dur %d", it.Duration())
	}
	if it.ResumeSeconds() != durResume {
		t.Fatalf("resume %d", it.ResumeSeconds())
	}
	z := Item{}
	if z.Duration() != 0 || z.ResumeSeconds() != 0 {
		t.Fatal("zero")
	}
}

func TestKind(t *testing.T) {
	if (Item{Type: TypeMovie}).Kind() != TypeMovie {
		t.Fatal("type")
	}
	if (Item{CollectionType: CollectionMovies}).Kind() != CollectionMovies {
		t.Fatal("collection")
	}
}
