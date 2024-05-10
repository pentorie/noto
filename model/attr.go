package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var StatusTypes = []string{"announce", "ongoing", "aired", "hiatus", "cancelled", "unknown"}
var Rating = []string{"g", "pg", "pg13", "r17", "rp", "rx"}
var TypeAnime = []string{"tv", "ova", "ona", "movie", "special", "music"}
var TypeRanobe = []string{"ln", "oneshot", "doujinshi"}
var TypeManga = []string{"manga", "oneshot", "manhwa", "manhua", "doujinshi"}
var Entity = []string{"anime", "manga", "ranobe", "news", "review", "user", "character", "person"}
var GenresList = []string{"action", "romance"}
var ModerationStatuses = []int{0, 1, 2}
var EntityTypes = []string{"anime", "manga", "ranobe"}
var EntitiesExt = []string{"anime", "manga", "ranobe", "character", "person", "news"}

// массив под жанры
type Genres []string
type Themes []string

func (g *Genres) Scan(value interface{}) error {
	if value == nil {
		*g = nil
		return nil
	}
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("genres: expected string value, got %T", value)
	}
	strValue = strings.Trim(strValue, "{}")           // remove curly braces
	strValue = strings.ReplaceAll(strValue, "\"", "") // remove double quotes
	*g = strings.Split(strValue, ",")
	return nil
}

func (g Genres) Value() (driver.Value, error) {
	if len(g) == 0 {
		return "{}", nil
	}

	return "{" + strings.Join(g, ",") + "}", nil
}

func (n *Themes) Scan(value interface{}) error {
	if value == nil {
		*n = nil
		return nil
	}
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("genres: expected string value, got %T", value)
	}
	strValue = strings.Trim(strValue, "{}")           // remove curly braces
	strValue = strings.ReplaceAll(strValue, "\"", "") // remove double quotes
	*n = strings.Split(strValue, ",")
	return nil
}

func (n Themes) Value() (driver.Value, error) {
	if len(n) == 0 {
		return "{}", nil
	}

	return "{" + strings.Join(n, ",") + "}", nil
}

type Altt struct {
	Eng  string `json:"eng"`
	Rus  string `json:"rus"`
	Orig string `json:"orig"`
	Etc  string `json:"etc"`
}

type MangaSources struct {
	Mangadex  string `json:"mangadex"`
	Readmanga string `json:"readmanga"`
}

type AnimeSources struct {
	Kodik string `json:"kodik"`
}

type RanobeSources struct {
	Ranobes string `json:"ranobes"`
}

type Marks struct {
	Mark1 int `json:"mark1"`
	Mark2 int `json:"mark2"`
	Mark3 int `json:"mark3"`
	Mark4 int `json:"mark4"`
	Mark5 int `json:"mark5"`
}

type Multilang struct {
	Eng  string `json:"eng"`
	Rus  string `json:"rus"`
	Orig string `json:"orig"`
	Etc  string `json:"etc"`
}

type RequestedList struct {
	TitleID  int       `json:"id"`
	Title    string    `json:"title"`
	Slug     string    `json:"slug"`
	Status   string    `json:"status"`
	Progress float32   `json:"progress"`
	Mark     int       `json:"mark"`
	Altt     Multilang `json:"altt"`
	Type     string    `json:"type"`
	Cover    string    `json:"cover"`
	Rating   string    `json:"agerating"`
	Genres   Genres    `json:"genres"`
	Themes   Genres    `json:"themes"`
	AiredOn  time.Time `json:"airedon"`
}

type ShortListQuery struct {
	Type   string `json:"type"`
	Cover  string `json:"cover"`
	Rating string `json:"agerating"`
	Title  string `json:"title"`
}

func (d Multilang) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Multilang) Scan(value interface{}) error {
	e, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(e, &d)
}

func (a Altt) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Altt) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (b Marks) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *Marks) Scan(value interface{}) error {
	c, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(c, &b)
}

func (e ExtLinks) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *ExtLinks) Scan(value interface{}) error {
	f, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(f, &e)
}
