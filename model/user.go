package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// type Favourites struct {
// 	Title_id  `json:"title_en"`
// 	Title_ru string `json:"title_ru"`
// 	Title_jp string `json:"title_jp"`
// }

type About struct {
	City   string `json:"city"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Link   string `json:"link"`
	Gender string `json:"gender"`
}

type User struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	CreatedAt   int64  `gorm:"autoCreateTime" json:"-" `
	UpdatedAt   int64  `gorm:"autoUpdateTime:milli" json:"-"`
	Username    string `json:"username"`
	Login       string `gorm:"unique" json:"login"`
	Email       string `json:"email" gorm:"unique"`
	Role        int    `gorm:"default:0" json:"role"`
	About       About  `gorm:"type:jsonb" json:"about"`
	Avatar      string `json:"avatar" gorm:"default:avatars/default.png"`
	Description string `json:"description"`
	Password    []byte `json:"-"`

	//News    []News    `gorm:"foreignKey:Author_id;references:ID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	//Review  []Review  `gorm:"foreignKey:Author_id;references:ID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	//Comment []Comment `gorm:"foreignKey:Author_id;references:ID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
}

type UserPublic struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	CreatedAt   int64  `gorm:"autoCreateTime" json:"-" `
	UpdatedAt   int64  `gorm:"autoUpdateTime:milli" json:"-"`
	Username    string `json:"username"`
	Login       string `gorm:"unique" json:"login"`
	Email       string `json:"email" gorm:"unique"`
	About       About  `gorm:"type:jsonb" json:"about"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}

func (a About) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *About) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
