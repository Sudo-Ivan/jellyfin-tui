package jellyfin

import "fmt"

type PublicInfo struct {
	ServerName string `json:"ServerName"`
	Version    string `json:"Version"`
	ID         string `json:"Id"`
}

type User struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

type AuthResult struct {
	AccessToken string `json:"AccessToken"`
	ServerID    string `json:"ServerId"`
	User        User   `json:"User"`
}

type QueryResult struct {
	Items            []Item `json:"Items"`
	TotalRecordCount int    `json:"TotalRecordCount"`
}

type Person struct {
	Name string `json:"Name"`
	Role string `json:"Role"`
	Type string `json:"Type"`
}

type UserData struct {
	PlaybackPositionTicks int64   `json:"PlaybackPositionTicks"`
	PlayCount             int     `json:"PlayCount"`
	IsFavorite            bool    `json:"IsFavorite"`
	Played                bool    `json:"Played"`
	PlayedPercentage      float64 `json:"PlayedPercentage"`
}

type Item struct {
	ID                 string   `json:"Id"`
	Name               string   `json:"Name"`
	Type               string   `json:"Type"`
	CollectionType     string   `json:"CollectionType"`
	Overview           string   `json:"Overview"`
	ProductionYear     int      `json:"ProductionYear"`
	DateCreated        string   `json:"DateCreated"`
	Genres             []string `json:"Genres"`
	People             []Person `json:"People"`
	RunTimeTicks       int64    `json:"RunTimeTicks"`
	IndexNumber        int      `json:"IndexNumber"`
	ParentIndexNumber  int      `json:"ParentIndexNumber"`
	ChildCount         int      `json:"ChildCount"`
	RecursiveItemCount int      `json:"RecursiveItemCount"`
	SeriesID           string   `json:"SeriesId"`
	SeriesName         string   `json:"SeriesName"`
	SeasonID           string   `json:"SeasonId"`
	SeasonName         string   `json:"SeasonName"`
	ParentID           string   `json:"ParentId"`
	UserData           UserData `json:"UserData"`
}

func (it Item) Duration() int64 {
	if it.RunTimeTicks <= 0 {
		return 0
	}
	return it.RunTimeTicks / ticksPerSec
}

func (it Item) ResumeSeconds() int {
	if it.UserData.PlaybackPositionTicks <= 0 {
		return 0
	}
	return int(it.UserData.PlaybackPositionTicks / ticksPerSec)
}

func (it Item) Label() string {
	switch it.Type {
	case TypeEpisode:
		s := it.SeriesName
		if s == "" {
			s = it.Name
		}
		if it.ParentIndexNumber > 0 && it.IndexNumber > 0 {
			return fmt.Sprintf("%s S%dE%d  %s", s, it.ParentIndexNumber, it.IndexNumber, it.Name)
		}
		if it.Name != "" && it.Name != s {
			return s + "  " + it.Name
		}
		return s
	case TypeMovie:
		if it.ProductionYear > 0 {
			return fmt.Sprintf("%s  (%d)", it.Name, it.ProductionYear)
		}
		return it.Name
	case TypeSeason:
		if it.IndexNumber > 0 {
			return fmt.Sprintf("Season %d", it.IndexNumber)
		}
		return it.Name
	default:
		return it.Name
	}
}

func (it Item) Kind() string {
	if it.Type != "" {
		return it.Type
	}
	return it.CollectionType
}
