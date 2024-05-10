package model

type Genre struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"-" `
	UpdatedAt  int64  `gorm:"autoUpdateTime" json:"-"`
	Genre_name string `gorm:"unique;type:varchar(50)" json:"genre_name"`
}
