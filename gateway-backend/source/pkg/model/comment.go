package model

// Комментарий к новостной публикации
type Comment struct {
	ID       int    `json:"ID"`
	PostId   int    `json:"PostId"`
	ParentId int    `json:"ParentId"`
	PubTime  int64  `json:"PubTime"`
	Content  string `json:"Content"`
}
