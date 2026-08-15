package jellyfin

import (
	"net/url"
	"strconv"
	"strings"
)

// ItemQuery is a GET /Items filter as defined by Jellyfin's ItemsController.
// Genres are pipe-delimited. Years are comma-delimited. Person is a single name.
type ItemQuery struct {
	UserID      string
	ParentID    string
	Search      string
	Include     string
	Fields      string
	SortBy      string
	SortOrd     string
	Limit       int
	Recurse     bool
	UserData    bool
	Genres      []string
	Years       []int
	Person      string
	PersonTypes []string
}

func (q ItemQuery) Values() url.Values {
	v := url.Values{}
	q.setCore(v)
	q.setFilters(v)
	return v
}

func (q ItemQuery) setCore(v url.Values) {
	setIf(v, qUserID, q.UserID)
	setIf(v, qParent, q.ParentID)
	setIf(v, qSearch, q.Search)
	setIf(v, qInclude, q.Include)
	setIf(v, qFields, q.Fields)
	setIf(v, qSortBy, q.SortBy)
	setIf(v, qSortOrd, q.SortOrd)
	if q.Limit > 0 {
		v.Set(qLimit, strconv.Itoa(q.Limit))
	}
	if q.Recurse {
		v.Set(qRecurse, boolTrue)
	} else {
		v.Set(qRecurse, boolFalse)
	}
	if q.UserData {
		v.Set(qUserData, boolTrue)
	}
}

func (q ItemQuery) setFilters(v url.Values) {
	if len(q.Genres) > 0 {
		v.Set(qGenres, strings.Join(q.Genres, genreSep))
	}
	if len(q.Years) > 0 {
		v.Set(qYears, joinInts(q.Years, yearSep))
	}
	setIf(v, qPerson, q.Person)
	if len(q.PersonTypes) > 0 {
		v.Set(qPersonTypes, strings.Join(q.PersonTypes, yearSep))
	}
}

func setIf(v url.Values, key, val string) {
	if val != "" {
		v.Set(key, val)
	}
}

func joinInts(ns []int, sep string) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, sep)
}
