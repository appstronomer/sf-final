package checker

// Комментарий к новостной публикации
type Comment struct {
	ID       int    `json:"ID"`
	PostId   int    `json:"PostId"`
	ParentId string `json:"ParentId"`
	PubTime  int64  `json:"PubTime"`
	Username string `json:"Username"`
	Content  string `json:"Content"`
}
