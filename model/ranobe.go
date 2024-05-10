package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Ranobe struct {
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
	Marks    Marks         `gorm:"type:jsonb" json:"marks"`
	MarkMean float32       `json:"markmean" gorm:"default:0.0001"`
	Sources  RanobeSources `gorm:"type:jsonb" json:"sources"`
	Duration int           `json:"duration" gorm:"default:24"`
	ExtLinks ExtLinks      `gorm:"type:jsonb" json:"extlinks"`
}

type RanobeWO struct {
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
	Sources  RanobeSources `gorm:"type:jsonb" json:"sources"`
	Duration int           `json:"duration" gorm:"default:24"`
	ExtLinks ExtLinks      `gorm:"type:jsonb" json:"extlinks"`
}

type UpdateRanobe struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	AuthorLog string `json:"author"`
	TitleID   int    `json:"titleid"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdat"`
	Changes   RanobeWO `json:"changes" gorm:"type:jsonb"`
	ReqStatus int    `json:"reqstatus" gorm:"default:0"`
}

type UpdateRanobeShort struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	AuthorLog string `json:"author"`
	TitleID   int    `json:"titleid"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdat"`
	ReqStatus int    `json:"reqstatus" gorm:"default:0"`
}

type RanobeTimings struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Status      string    `json:"status"`
	AiredOn     time.Time `json:"airedon"`
	Episodes    int       `json:"episodes"`
	CurrEpisode int       `json:"currepisodes" gorm:"default:0"`
	Period      int       `json:"period" gorm:"default:7"`
	NextEpisode time.Time `json:"nextepisode"`
}

type RanobeShort struct {
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Slug     string   `gorm:"unique" json:"slug"`
	Altt    Multilang  `gorm:"type:jsonb" json:"altt"`
	Cover    string   `json:"cover"`
}

func (x Ranobe) Value() (driver.Value, error) {
	return json.Marshal(x)
}

func (x *Ranobe) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &x)
}

func (ad RanobeWO) Value() (driver.Value, error) {
	return json.Marshal(ad)
}

func (ad *RanobeWO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ad)
}

func (ab RanobeSources) Value() (driver.Value, error) {
	return json.Marshal(ab)
}

func (ab *RanobeSources) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ab)
}

// преобразования на единственное число для таблицы
func (Ranobe) TableName() string {
	return "ranobe"
}

func (RanobeShort) TableName() string {
	return "ranobe"
}
