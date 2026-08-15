package jellyfin

import (
	"errors"
	"net/url"
)

var errNoMovie = errors.New("no movie matched")

func (c *Client) items(q ItemQuery) ([]Item, error) {
	if q.UserID == "" {
		q.UserID = c.UserID
	}
	var r QueryResult
	err := c.get(pathItems, q.Values(), &r)
	return r.Items, err
}

// Views lists the user's libraries.
func (c *Client) Views() ([]Item, error) {
	q := url.Values{qUserID: {c.UserID}}
	var r QueryResult
	err := c.get(pathUserViews, q, &r)
	return r.Items, err
}

// Resume lists in-progress items.
func (c *Client) Resume(limit int) ([]Item, error) {
	q := url.Values{
		qUserID: {c.UserID},
		qLimit:  {limitQ(limit)},
	}
	var r QueryResult
	err := c.get(pathResume, q, &r)
	return r.Items, err
}

// NextUp lists the next episode for in-progress series.
func (c *Client) NextUp(limit int) ([]Item, error) {
	q := url.Values{
		qUserID: {c.UserID},
		qLimit:  {limitQ(limit)},
		qFields: {fieldsHome},
	}
	var r QueryResult
	err := c.get(pathNextUp, q, &r)
	return r.Items, err
}

// NewlyAdded lists items sorted by DateCreated descending.
func (c *Client) NewlyAdded(limit int) ([]Item, error) {
	return c.items(ItemQuery{
		Include:  includeLatest,
		Fields:   fieldsLatest,
		SortBy:   sortCreated,
		SortOrd:  orderDesc,
		Limit:    limit,
		Recurse:  true,
		UserData: true,
	})
}

// Latest is NewlyAdded. Kept for callers that used the old name.
func (c *Client) Latest(limit int) ([]Item, error) {
	return c.NewlyAdded(limit)
}

// Children lists items under parentID. includeTypes is a comma list or empty.
func (c *Client) Children(parentID, includeTypes string) ([]Item, error) {
	q := ItemQuery{
		ParentID: parentID,
		Fields:   fieldsChildren,
		SortBy:   sortNameIndex,
		SortOrd:  orderAsc,
		UserData: true,
	}
	if includeTypes != "" {
		q.Include = includeTypes
		q.Recurse = true
	}
	return c.items(q)
}

// Search finds movies, series, and episodes by name and optional filters.
func (c *Client) Search(term string, limit int) ([]Item, error) {
	return c.Query(ParseFilter(term).Query(limit))
}

// Query runs GET /Items with q.
func (c *Client) Query(q ItemQuery) ([]Item, error) {
	if q.Fields == "" {
		q.Fields = fieldsSearch
	}
	if q.Include == "" {
		q.Include = includeSearch
	}
	return c.items(q)
}

// RandomMovie returns one movie using SortBy=Random.
func (c *Client) RandomMovie(parentID string) (Item, error) {
	list, err := c.items(ItemQuery{
		ParentID: parentID,
		Include:  includeMovies,
		Fields:   fieldsItem,
		SortBy:   sortRandom,
		Limit:    1,
		Recurse:  true,
		UserData: true,
	})
	if err != nil {
		return Item{}, err
	}
	if len(list) == 0 {
		return Item{}, errNoMovie
	}
	return list[0], nil
}

// Item loads one library object by id.
func (c *Client) Item(id string) (Item, error) {
	var it Item
	q := url.Values{
		qUserID: {c.UserID},
		qFields: {fieldsItem},
	}
	err := c.get(pathItems+"/"+id, q, &it)
	return it, err
}

// Seasons lists seasons for a series.
func (c *Client) Seasons(seriesID string) ([]Item, error) {
	q := url.Values{
		qUserID: {c.UserID},
		qFields: {fieldsSeason},
	}
	var r QueryResult
	err := c.get("/Shows/"+seriesID+"/Seasons", q, &r)
	return r.Items, err
}

// Episodes lists episodes in a series, optionally limited to seasonID.
func (c *Client) Episodes(seriesID, seasonID string) ([]Item, error) {
	q := url.Values{
		qUserID: {c.UserID},
		qFields: {fieldsEpisode},
	}
	if seasonID != "" {
		q.Set(qSeason, seasonID)
	}
	var r QueryResult
	err := c.get("/Shows/"+seriesID+"/Episodes", q, &r)
	return r.Items, err
}

// AdjacentEpisodes returns the current episode plus neighbors.
func (c *Client) AdjacentEpisodes(seriesID, episodeID string) ([]Item, error) {
	q := url.Values{
		qUserID:   {c.UserID},
		qAdjacent: {episodeID},
		qFields:   {fieldsAdjacent},
	}
	var r QueryResult
	err := c.get("/Shows/"+seriesID+"/Episodes", q, &r)
	return r.Items, err
}
