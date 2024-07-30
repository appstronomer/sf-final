package api

import (
	"sf-news/pkg/storage"
)

type PostsCollection struct {
	Data      []storage.PostShort `json:"Data"`
	ElemCount int                 `json:"ElemCount"`
	PageCount int                 `json:"PageCount"`
	PageNo    int                 `json:"PageNo"`
}
