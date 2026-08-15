package jellyfin

import "testing"

func TestItemMatchesFilter(t *testing.T) {
	it := Item{
		Name:           "Blade Runner",
		SeriesName:     "",
		ProductionYear: yr1982,
		Genres:         []string{"Sci-Fi", "Thriller"},
		People: []Person{
			{Name: "Harrison Ford", Type: "Actor"},
			{Name: "Ridley Scott", Type: "Director"},
		},
	}
	if !it.Matches(Filter{}) {
		t.Fatal("empty filter")
	}
	if !it.Matches(Filter{Term: nameBlade}) {
		t.Fatal("term")
	}
	if it.Matches(Filter{Term: "matrix"}) {
		t.Fatal("term miss")
	}
	if !it.Matches(Filter{Genres: []string{"sci-fi"}}) {
		t.Fatal("genre fold")
	}
	if !it.Matches(Filter{Genres: []string{"Sci-Fi", "Thriller"}}) {
		t.Fatal("genre and")
	}
	if it.Matches(Filter{Genres: []string{"Comedy"}}) {
		t.Fatal("genre miss")
	}
	if !it.Matches(Filter{Person: "ford"}) {
		t.Fatal("actor")
	}
	if it.Matches(Filter{Person: "scott"}) {
		t.Fatal("director must not match actor filter")
	}
	if !it.Matches(Filter{Years: []int{yr1982}}) {
		t.Fatal("year")
	}
	if it.Matches(Filter{Years: []int{yr1999}}) {
		t.Fatal("year miss")
	}
}

func TestItemMatchesEpisodeTermOnSeriesName(t *testing.T) {
	it := Item{Name: "Pilot", SeriesName: "The Wire", Type: TypeEpisode}
	if !it.Matches(Filter{Term: "wire"}) {
		t.Fatal("series name")
	}
}

func TestMatchPersonEmptyType(t *testing.T) {
	it := Item{People: []Person{{Name: nameFord}}}
	if !it.Matches(Filter{Person: nameFord}) {
		t.Fatal("untyped person")
	}
}
