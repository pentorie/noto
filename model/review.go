package model

type Review struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"createdat" `
	UpdatedAt  int64  `gorm:"autoUpdateTime:milli" json:"-"`
	Title_id   int    `json:"title_id"`
	Title_type string `json:"title_type"`
	Author_id  int    `json:"author_id"`
	Author_log string `json:"author_log"`
	Content    string `json:"content"`
	Rated      int    `json:"rated"`
}

type ReviewExp struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"createdat" `
	UpdatedAt  int64  `gorm:"autoUpdateTime:milli" json:"-"`
	Title_id   int    `json:"title_id"`
	Title_type string `json:"title_type"`
	Author_id  int    `json:"author_id"`
	AuthUname  string `json:"authuname"`
	AuthLog    string `json:"authlog"`
	AuthAvatar string `json:"authavatar"`
	Content    string `json:"content"`
	Rated      int    `json:"rated"`
}
