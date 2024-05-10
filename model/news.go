package model

type News struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"createdat"`
	UpdatedAt  int64  `gorm:"autoUpdateTime" json:"-"`
	Author_id  int    `json:"authorid"`
	AuthLog    string `json:"authorlog"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Cover      string `json:"cover"`
	CommentQty int    `json:"commentqty" gorm:"default:0"`
}

type NewsShort struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"createdat"`
	AuthLog    string `json:"authorlog"`
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	CommentQty int    `json:"commentqty" gorm:"default:0"`
}
