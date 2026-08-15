package jellyfin

func NextEpisode(list []Item, currentID string) (Item, bool) {
	idx := indexOf(list, currentID)
	if idx >= 0 && idx+1 < len(list) {
		return list[idx+1], true
	}
	cur, ok := findItem(list, currentID)
	if !ok {
		return Item{}, false
	}
	return nextBySeason(list, cur)
}

func indexOf(list []Item, id string) int {
	for i, it := range list {
		if it.ID == id {
			return i
		}
	}
	return -1
}

func findItem(list []Item, id string) (Item, bool) {
	for _, it := range list {
		if it.ID == id {
			return it, true
		}
	}
	return Item{}, false
}

func nextBySeason(list []Item, cur Item) (Item, bool) {
	best := Item{}
	found := false
	for _, it := range list {
		if it.ID == cur.ID || !seasonAhead(it, cur) {
			continue
		}
		if !found || episodeLess(it, best) {
			best = it
			found = true
		}
	}
	return best, found
}

func seasonAhead(a, cur Item) bool {
	if a.ParentIndexNumber > cur.ParentIndexNumber {
		return true
	}
	return a.ParentIndexNumber == cur.ParentIndexNumber && a.IndexNumber > cur.IndexNumber
}

func episodeLess(a, b Item) bool {
	if a.ParentIndexNumber != b.ParentIndexNumber {
		return a.ParentIndexNumber < b.ParentIndexNumber
	}
	return a.IndexNumber < b.IndexNumber
}
