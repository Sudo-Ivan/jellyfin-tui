package jellyfin

import (
	"net/url"
	"testing"
)

func TestItemQueryEncodesJellyfinBinders(t *testing.T) {
	q := ItemQuery{
		UserID:      "user-1",
		ParentID:    "lib-9",
		Search:      nameBlade,
		Include:     includeMovies,
		Fields:      fieldsSearch,
		SortBy:      sortRandom,
		Limit:       1,
		Recurse:     true,
		UserData:    true,
		Genres:      []string{nameAction, "Thriller"},
		Years:       []int{yr1999, yr2000},
		Person:      nameFord,
		PersonTypes: []string{personActor},
	}
	got := q.Values()
	want := url.Values{}
	want.Set("UserId", "user-1")
	want.Set("ParentId", "lib-9")
	want.Set("SearchTerm", nameBlade)
	want.Set("IncludeItemTypes", "Movie")
	want.Set("Fields", fieldsSearch)
	want.Set("SortBy", "Random")
	want.Set("Limit", "1")
	want.Set("Recursive", "true")
	want.Set("EnableUserData", "true")
	want.Set("Genres", nameAction+"|Thriller")
	want.Set("Years", "1999,2000")
	want.Set("Person", nameFord)
	want.Set("PersonTypes", "Actor")
	if got.Encode() != want.Encode() {
		t.Fatalf("query encode\n got %s\nwant %s", got.Encode(), want.Encode())
	}
}

func TestNewlyAddedQuerySort(t *testing.T) {
	q := ItemQuery{
		Include: includeLatest,
		SortBy:  sortCreated,
		SortOrd: orderDesc,
		Limit:   limitAdded,
		Recurse: true,
	}
	v := q.Values()
	if v.Get("SortBy") != "DateCreated" {
		t.Fatalf("SortBy=%q", v.Get("SortBy"))
	}
	if v.Get("SortOrder") != "Descending" {
		t.Fatalf("SortOrder=%q", v.Get("SortOrder"))
	}
	if v.Get("IncludeItemTypes") != "Movie,Episode,Series" {
		t.Fatalf("include=%q", v.Get("IncludeItemTypes"))
	}
}

func TestFilterQuerySetsActorType(t *testing.T) {
	q := Filter{Person: "Hanks", Term: "castaway"}.Query(limitSearch)
	if q.PersonTypes[0] != personActor {
		t.Fatal("personTypes")
	}
	if q.Search != "castaway" {
		t.Fatal("term")
	}
	v := q.Values()
	if v.Get("Person") != "Hanks" || v.Get("PersonTypes") != "Actor" {
		t.Fatalf("%v", v)
	}
}

func TestEmptyOptionalFiltersOmitted(t *testing.T) {
	v := ItemQuery{UserID: "u", Recurse: false}.Values()
	if v.Get("Genres") != "" || v.Get("Years") != "" || v.Get("Person") != "" {
		t.Fatalf("empty filters leaked: %v", v)
	}
	if v.Get("Recursive") != "false" {
		t.Fatal("folder listing must send Recursive=false")
	}
}
