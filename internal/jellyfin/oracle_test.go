package jellyfin

import (
	"fmt"
	"net/url"
	"testing"
)

func TestOracleJellyfinItemsContract(t *testing.T) {
	oracleGenres(t)
	oracleYears(t)
	oracleParseQuery(t)
	oracleRandom(t)
	oracleNewlyAdded(t)
	oraclePaths(t)
	oracleIndependentEncode(t)
	_, _ = fmt.Println("QUERY_ORACLE_PROVED")
}

func oracleGenres(t *testing.T) {
	t.Helper()
	got := ItemQuery{Genres: []string{nameAction, "Sci-Fi"}}.Values().Get("Genres")
	if got != nameAction+"|Sci-Fi" {
		t.Fatalf("genres binder want pipe got %q", got)
	}
}

func oracleYears(t *testing.T) {
	t.Helper()
	got := ItemQuery{Years: []int{yr1999, yr2001}}.Values().Get("Years")
	if got != "1999,2001" {
		t.Fatalf("years binder want comma got %q", got)
	}
}

func oracleParseQuery(t *testing.T) {
	t.Helper()
	f := ParseFilter(`genre:Action actor:"Harrison Ford" year:1982 blade`)
	q := f.Query(durResume).Values()
	if q.Get("Genres") != nameAction || q.Get("Person") != "Harrison Ford" || q.Get("PersonTypes") != "Actor" {
		t.Fatalf("parse->query %s", q.Encode())
	}
	if q.Get("Years") != "1982" || q.Get("SearchTerm") != nameBlade {
		t.Fatalf("term/year %s", q.Encode())
	}
}

func oracleRandom(t *testing.T) {
	t.Helper()
	rnd := ItemQuery{SortBy: sortRandom, Include: includeMovies, Limit: 1, Recurse: true}.Values()
	if rnd.Get("SortBy") != "Random" || rnd.Get("IncludeItemTypes") != "Movie" || rnd.Get("Limit") != "1" {
		t.Fatalf("random %s", rnd.Encode())
	}
}

func oracleNewlyAdded(t *testing.T) {
	t.Helper()
	added := ItemQuery{SortBy: sortCreated, SortOrd: orderDesc, Include: includeLatest}.Values()
	if added.Get("SortBy") != "DateCreated" || added.Get("SortOrder") != "Descending" {
		t.Fatalf("newly added %s", added.Encode())
	}
}

func oraclePaths(t *testing.T) {
	t.Helper()
	if pathItems != "/Items" || pathResume != "/UserItems/Resume" || pathUserViews != "/UserViews" {
		t.Fatalf("paths items=%s resume=%s views=%s", pathItems, pathResume, pathUserViews)
	}
}

func oracleIndependentEncode(t *testing.T) {
	t.Helper()
	u := url.Values{}
	u.Set("Genres", nameAction+"|Sci-Fi")
	u.Set("Years", "1999,2001")
	u.Set("Recursive", "false")
	got := ItemQuery{Genres: []string{nameAction, "Sci-Fi"}, Years: []int{yr1999, yr2001}}.Values().Encode()
	if got != u.Encode() {
		t.Fatal("independent url.Values oracle mismatch")
	}
}
