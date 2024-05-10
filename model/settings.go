package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Settings struct {
	UserID      int         `json:"userid"`
	ColorScheme string      `json:"colorscheme"`
	DataCache   bool        `json:"datacache"`
	Language    LanguageSet `json:"language" gorm:"type:jsonb"`
}

// предпочительные языки
type LanguageSet struct {
	InterfaceLanguage string `json:"interface"`
	GenresLanguage    string `json:"genres"`
	PrefMainTitle     string `json:"maintitle"`
	PrefSubTitle      string `json:"subtitle"`
	PrefDescription   string `json:"description"`
}

func (a LanguageSet) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *LanguageSet) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (b Settings) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *Settings) Scan(value interface{}) error {
	c, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(c, &b)
}
