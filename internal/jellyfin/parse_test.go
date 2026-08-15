package jellyfin

import (
	"reflect"
	"testing"
)

func TestParseFilterTable(t *testing.T) {
	cases := []struct {
		in   string
		want Filter
	}{
		{nameBlade, Filter{Term: nameBlade}},
		{"genre:" + nameAction, Filter{Genres: []string{nameAction}}},
		{"g:Sci-Fi g:" + nameAction, Filter{Genres: []string{"Sci-Fi", nameAction}}},
		{"genre:" + nameAction + "|Thriller", Filter{Genres: []string{nameAction, "Thriller"}}},
		{"actor:" + nameFord, Filter{Person: nameFord}},
		{`actor:"Tom Hanks"`, Filter{Person: "Tom Hanks"}},
		{"person:" + nameFord + " " + nameBlade, Filter{Term: nameBlade, Person: nameFord}},
		{"year:1999", Filter{Years: []int{yr1999}}},
		{"y:2010-2012", Filter{Years: []int{yr2010, yr2011, yr2012}}},
		{"year:2012-2010", Filter{Years: []int{yr2010, yr2011, yr2012}}},
		{"year:2000,2002", Filter{Years: []int{yr2000, yr2002}}},
		{"genre:" + nameAction + " actor:" + nameFord + " year:1999 " + nameBlade, Filter{
			Term: nameBlade, Genres: []string{nameAction}, Person: nameFord, Years: []int{yr1999},
		}},
		{"year:12", Filter{}},
		{"actor:" + nameFord + " actor:Hanks", Filter{Person: nameFord}},
	}
	for _, tc := range cases {
		got := ParseFilter(tc.in)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%q\n got %#v\nwant %#v", tc.in, got, tc.want)
		}
	}
}

func TestParseFilterEmpty(t *testing.T) {
	if !ParseFilter("").Empty() || !ParseFilter("   ").Empty() {
		t.Fatal("blank")
	}
}

func TestYearRangeCap(t *testing.T) {
	ys := parseYears("1888-2100")
	if len(ys) != yearSpanCap+1 {
		t.Fatalf("len=%d", len(ys))
	}
	if ys[0] != yearMin || ys[len(ys)-1] != yearMin+yearSpanCap {
		t.Fatalf("%v", ys[0])
	}
}
