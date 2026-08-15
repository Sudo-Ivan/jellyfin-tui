package jellyfin

import (
	"strconv"
	"strings"
	"unicode"
)

// Filter is a user-typed search: term plus optional genre, actor, and year constraints.
type Filter struct {
	Term   string
	Genres []string
	Person string
	Years  []int
}

func (f Filter) Empty() bool {
	return f.Term == "" && len(f.Genres) == 0 && f.Person == "" && len(f.Years) == 0
}

func (f Filter) Query(limit int) ItemQuery {
	q := ItemQuery{
		Search:   f.Term,
		Include:  includeSearch,
		Fields:   fieldsSearch,
		Limit:    limit,
		Recurse:  true,
		UserData: true,
		Genres:   f.Genres,
		Years:    f.Years,
		Person:   f.Person,
	}
	if f.Person != "" {
		q.PersonTypes = []string{personActor}
	}
	return q
}

// ParseFilter splits a search line into term and genre:/actor:/year: fields.
// Quotes keep multi-word values together. Year ranges like 2010-2012 expand.
func ParseFilter(s string) Filter {
	var f Filter
	var terms []string
	for _, tok := range quotedFields(s) {
		key, val, ok := strings.Cut(tok, ":")
		if !ok || val == "" || !filterKey(key) {
			terms = append(terms, tok)
			continue
		}
		applyToken(&f, strings.ToLower(key), val)
	}
	f.Term = strings.TrimSpace(strings.Join(terms, " "))
	return f
}

func applyToken(f *Filter, key, val string) {
	switch key {
	case "genre", "g":
		f.Genres = append(f.Genres, splitList(val, genreSep)...)
	case "actor", "person", "p":
		if f.Person == "" {
			f.Person = val
		}
	case "year", "y":
		f.Years = append(f.Years, parseYears(val)...)
	}
}

func filterKey(key string) bool {
	switch strings.ToLower(key) {
	case "genre", "g", "actor", "person", "p", "year", "y":
		return true
	default:
		return false
	}
}

func splitList(s, sep string) []string {
	var out []string
	for p := range strings.SplitSeq(s, sep) {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseYears(s string) []int {
	s = strings.TrimSpace(s)
	if a, b, ok := strings.Cut(s, "-"); ok {
		return yearRange(atoiYear(a), atoiYear(b))
	}
	var out []int
	for p := range strings.SplitSeq(s, yearSep) {
		if y := atoiYear(p); y != 0 {
			out = append(out, y)
		}
	}
	return out
}

func yearRange(a, b int) []int {
	if a == 0 || b == 0 {
		return nil
	}
	if a > b {
		a, b = b, a
	}
	if b-a > yearSpanCap {
		b = a + yearSpanCap
	}
	out := make([]int, 0, b-a+1)
	for y := a; y <= b; y++ {
		out = append(out, y)
	}
	return out
}

func atoiYear(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n < yearMin || n > yearMax {
		return 0
	}
	return n
}

func quotedFields(s string) []string {
	var out []string
	var b strings.Builder
	inQ := false
	flush := func() {
		if b.Len() == 0 {
			return
		}
		out = append(out, b.String())
		b.Reset()
	}
	for _, r := range s {
		switch {
		case r == '"':
			inQ = !inQ
		case unicode.IsSpace(r) && !inQ:
			flush()
		default:
			_, _ = b.WriteRune(r)
		}
	}
	flush()
	return out
}
