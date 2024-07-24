package storage

// Новостная публикация с контентом
type PostFull struct {
	ID      int    `json:"ID"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

// Новостная публикация без контента
type PostShort struct {
	ID      int    `json:"ID"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
	Title   string `json:"Title"`
}

// Interface задаёт контракт на работу с БД.
type StorageIface interface {
	// Получение деталлизированной новости
	GetPost(id int) (PostFull, error)

	// Получение списка самых новых публикаций из хранилища новостей
	GetPosts(offset, limit int) ([]PostShort, error)

	// Получение списка самых новых публикаций из хранилища новостей
	// с учетом поиска
	FindPosts(search string, offset, limit int) ([]PostShort, error)

	// Получение количества новостей всего
	GetCount() (int, error)

	// Получение количества новостей с учетом поиска
	FindCount(search string) (int, error)

	// Добавление публикаций в хранилище новостей
	PushPosts([]PostFull) error
}
