package jellyfin

import "slices"

import "strings"

// Matches reports whether it satisfies f. Empty constraints are ignored.
func (it Item) Matches(f Filter) bool {
	if f.Empty() {
		return true
	}
	if f.Term != "" && !containsFold(it.Name, f.Term) && !containsFold(it.SeriesName, f.Term) {
		return false
	}
	if len(f.Genres) > 0 && !matchGenres(it.Genres, f.Genres) {
		return false
	}
	if f.Person != "" && !matchPerson(it.People, f.Person) {
		return false
	}
	if len(f.Years) > 0 && !matchYear(it.ProductionYear, f.Years) {
		return false
	}
	return true
}

func matchGenres(have, want []string) bool {
	for _, w := range want {
		ok := false
		for _, h := range have {
			if containsFold(h, w) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func matchPerson(people []Person, want string) bool {
	for _, p := range people {
		if p.Type != "" && !strings.EqualFold(p.Type, personActor) {
			continue
		}
		if containsFold(p.Name, want) {
			return true
		}
	}
	return false
}

func matchYear(y int, years []int) bool {
	return slices.Contains(years, y)
}

func containsFold(s, q string) bool {
	if q == "" {
		return true
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(q))
}
