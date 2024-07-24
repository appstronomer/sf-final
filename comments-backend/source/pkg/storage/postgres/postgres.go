package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"sf-comments/pkg/storage"
)

// Хранилище данных
type Storage struct {
	dbPool *pgxpool.Pool
}

// Конструктор объекта хранилища
func New(constr string) (*Storage, error) {
	dbPool, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	return &Storage{dbPool: dbPool}, nil
}

// Получает список комментариев
func (s *Storage) GetComments(postId, parentId, lastId, limit int) ([]storage.Comment, error) {
	var parentIdSql interface{}
	if parentId == 0 {
		parentIdSql = nil
	} else {
		parentIdSql = parentId
	}
	rows, err := s.dbPool.Query(context.Background(),
		`SELECT c.id, c.post_id, c.parent_id, c.pub_time, c.content
		FROM comments c
		WHERE c.post_id = $1 AND c.parent_id = $2 AND c.id > $3
		ORDER BY c.id ASC
		LIMIT $4;`,
		postId, parentIdSql, lastId,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []storage.Comment
	for rows.Next() {
		var c storage.Comment
		err := rows.Scan(&c.ID, &c.PostId, &c.ParentId, &c.PubTime, c.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// Добавляет комментарий к новости или другому комментарию
func (s *Storage) PushComment(comment storage.Comment) error {
	_, err := s.dbPool.Exec(context.Background(),
		`INSERT INTO posts (post_id, parent_id, pub_time, content)
		VALUES ($1, $2, $3, $4);`,
		comment.PostId,
		comment.ParentId,
		comment.PubTime,
		comment.Content,
	)
	return err
}
