package storage

// Комментарий к новостной публикации
type Comment struct {
	ID       int    `json:"ID"`
	PostId   int    `json:"PostId"`
	ParentId string `json:"ParentId"`
	PubTime  int64  `json:"PubTime"`
	Username string `json:"Username"`
	Content  string `json:"Content"`
}

// Interface задаёт контракт на работу с БД.
type StorageIface interface {

	// Получение списка комментариев к статье в режиме потока
	// postId - id новостной публикации
	// parentId - id родительского комментария (0 для корневых комментариев)
	// lastId - id последнего полученного потребителем комментария (0, если
	// это первый запрос комментариев)
	// limit - максимальное количество комментариев
	GetComments(postId, parentId, lastId, limit int) ([]Comment, error)

	// Добавление комментария в хранилище
	PushComment(Comment) error
}
