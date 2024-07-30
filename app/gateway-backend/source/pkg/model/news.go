package model

// Новостная публикация с контентом
type NewsFullDetailed struct {
	ID      int    `json:"ID"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

// Новостная публикация без контента
type NewsShortDetailed struct {
	ID      int    `json:"ID"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
	Title   string `json:"Title"`
}

// Коллекция новостных публицкаций с пагенцией
type NewsCollection struct {
	Data      []NewsShortDetailed `json:"Data"`
	ElemCount int                 `json:"ElemCount"`
	PageCount int                 `json:"PageCount"`
	PageNo    int                 `json:"PageNo"`
}

type NewsComplex struct {
	Data     NewsFullDetailed `json:"Data"`
	Comments []Comment        `json:"Comments"`
}
