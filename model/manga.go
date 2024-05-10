package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Manga struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"-" `
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"-"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
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
	//etc
	Marks    Marks        `gorm:"type:jsonb" json:"marks"`
	MarkMean float32      `json:"markmean" gorm:"default:0.0001"`
	Sources  MangaSources `gorm:"type:jsonb" json:"sources"`
	Duration int          `json:"duration" gorm:"default:24"`
	ExtLinks ExtLinks     `gorm:"type:jsonb" json:"extlinks"`
}

type MangaWO struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"-"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
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
	//etc
	Sources  MangaSources `gorm:"type:jsonb" json:"sources"`
	Duration int          `json:"duration" gorm:"default:24"`
	ExtLinks ExtLinks     `gorm:"type:jsonb" json:"extlinks"`
}

type UpdateManga struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	AuthorLog string `json:"author"`
	TitleID   int    `json:"titleid"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdat"`
	Changes   MangaWO  `json:"changes" gorm:"type:jsonb"`
	ReqStatus int    `json:"reqstatus" gorm:"default:0"`
}

type UpdateMangaShort struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	AuthorLog string `json:"author"`
	TitleID   int    `json:"titleid"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdat"`
	ReqStatus int    `json:"reqstatus" gorm:"default:0"`
}

type MangaTimings struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Status      string    `json:"status"`
	AiredOn     time.Time `json:"airedon"`
	Episodes    int       `json:"episodes"`
	CurrEpisode int       `json:"currepisodes" gorm:"default:0"`
	Period      int       `json:"period" gorm:"default:7"`
	NextEpisode time.Time `json:"nextepisode"`
}

type MangaShort struct {
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Slug     string   `gorm:"unique" json:"slug"`
	Altt    Multilang  `gorm:"type:jsonb" json:"altt"`
	Cover    string   `json:"cover"`
}

func (z Manga) Value() (driver.Value, error) {
	return json.Marshal(z)
}

func (z *Manga) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &z)
}

func (ac MangaWO) Value() (driver.Value, error) {
	return json.Marshal(ac)
}

func (ac *MangaWO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ac)
}

func (aa MangaSources) Value() (driver.Value, error) {
	return json.Marshal(aa)
}

func (aa *MangaSources) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &aa)
}

// преобразования на единственное число для таблицы
func (Manga) TableName() string {
	return "manga"
}

func (MangaShort) TableName() string {
	return "manga"
}
