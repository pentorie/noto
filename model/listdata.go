package model

type Listdata struct {
	TitleID  int     `json:"titleid"`
	Status   string  `json:"status"`
	Progress float32 `json:"progress"`
	Mark     int     `json:"mark"`
}
