package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Anime struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"-" `
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"-"`
	Title       string    `json:"title"`
	Studio      string    `json:"studio"`
	Slug        string    `gorm:"unique" json:"slug"`
	Altt    Multilang  `gorm:"type:jsonb" json:"altt"`
	Description Multilang `json:"description" gorm:"type:jsonb"`
	Type        string    `json:"type"`
	Cover       string    `json:"cover"`
	Status      string    `json:"status"`
	Rating      string    `json:"rating"`
	Genres      Genres    `json:"genres"`
	Themes      Genres    `json:"themes"`
	//расчёт для онгоингов
	AiredOn     time.Time `json:"airedon"`
	AiredEnd    time.Time `json:"airedend"`
	Episodes    int       `json:"episodes"`
	CurrEpisode int       `json:"currepisodes" gorm:"default:0"`
	Period      int       `json:"period" gorm:"default:7"`
	NextEpisode time.Time `json:"nextepisode"`
	//etc
	Marks    Marks        `gorm:"type:jsonb" json:"marks"`
	MarkMean float32      `json:"markmean" gorm:"default:0.0001"`
	Sources  AnimeSources `gorm:"type:jsonb" json:"sources"`
	Duration int          `json:"duration" gorm:"default:24"`
	ExtLinks ExtLinks     `gorm:"type:jsonb" json:"extlinks"`
}

type AnimeWO struct {
    ID          int       `gorm:"primaryKey" json:"id"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"-"`
	Title       string    `json:"title"`
	Studio      string    `json:"studio"`
	Slug        string    `gorm:"unique" json:"slug"`
	Altt    Multilang  `gorm:"type:jsonb" json:"altt"`
	Description Multilang `json:"description" gorm:"type:jsonb"`
	Type        string    `json:"type"`
	Cover       string    `json:"cover"`
	Status      string    `json:"status"`
	Rating      string    `json:"rating"`
	Genres      Genres    `json:"genres"`
	Themes      Genres    `json:"themes"`
	//расчёт для онгоингов
	AiredOn     time.Time `json:"airedon"`
	AiredEnd    time.Time `json:"airedend"`
	Episodes    int       `json:"episodes"`
	CurrEpisode int       `json:"currepisodes" gorm:"default:0"`
	Period      int       `json:"period" gorm:"default:7"`
	NextEpisode time.Time `json:"nextepisode"`
	//etc
	Sources  AnimeSources `gorm:"type:jsonb" json:"sources"`
	Duration int          `json:"duration" gorm:"default:24"`
	ExtLinks ExtLinks     `gorm:"type:jsonb" json:"extlinks"`
}

type UpdateAnime struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	AuthorLog string `json:"author"`
	TitleID   int    `json:"titleid"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdat"`
	Changes   AnimeWO  `json:"changes" gorm:"type:jsonb"`
	ReqStatus int    `json:"reqstatus" gorm:"default:0"`
}

type UpdateAnimeShort struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	AuthorLog string `json:"author"`
	TitleID   int    `json:"titleid"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdat"`
	ReqStatus int    `json:"reqstatus" gorm:"default:0"`
}

type ExtLinks struct {
	Official  string `json:"official"`
	Youtube   string `json:"youtube"`
	Twitter   string `json:"twitter"`
	Shikimori string `json:"shikimori"`
	MAL       string `json:"mal"`
	Kitsu     string `json:"kitsu"`
	AniList   string `json:"anilist"`
}

type AnimeTimings struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Status      string    `json:"status"`
	AiredOn     time.Time `json:"airedon"`
	Episodes    int       `json:"episodes"`
	CurrEpisode int       `json:"currepisodes" gorm:"default:0"`
	Period      int       `json:"period" gorm:"default:7"`
	NextEpisode time.Time `json:"nextepisode"`
}

type AnimeShort struct {
	Title    string   `json:"title"`
	Studio   string   `json:"studio"`
	Slug     string   `gorm:"unique" json:"slug"`
	Altt    Multilang  `gorm:"type:jsonb" json:"altt"`
	Cover    string   `json:"cover"`
}

// велью-сканнеры для json-вложенностей в модели
func (y Anime) Value() (driver.Value, error) {
	return json.Marshal(y)
}

func (y *Anime) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &y)
}

func (r AnimeWO) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *AnimeWO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &r)
}

func (c AnimeSources) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *AnimeSources) Scan(value interface{}) error {
	d, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(d, &c)
}

// преобразования на единственное число для таблицы
func (Anime) TableName() string {
	return "anime"
}

func (AnimeShort) TableName() string {
	return "anime"
}
