package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"sf-news/pkg/storage"
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

// Возвращает новость по её id
func (s *Storage) GetPost(id int) (storage.PostFull, error) {
	var p storage.PostFull
	row := s.dbPool.QueryRow(context.Background(),
		`SELECT p.id, p.pub_time, p.link, p.title, p.content
		FROM posts p
		WHERE p.id = $1;`,
		id,
	)
	err := row.Scan(&p.ID, &p.PubTime, &p.Link, &p.Title, &p.Content)
	if err != nil {
		return storage.PostFull{}, err
	}
	return p, nil
}

// Добавляет новости в хранилище
func (s *Storage) PushPosts(posts []storage.PostFull) error {
	for _, post := range posts {
		_, err := s.dbPool.Exec(context.Background(), `
			INSERT INTO posts (pub_time, link, title, content)
			VALUES ($1, $2, $3, $4);`,
			post.PubTime,
			post.Link,
			post.Title,
			post.Content,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Получает самые свежие новости из хранилища
func (s *Storage) GetPosts(offset, limit int) ([]storage.PostShort, error) {
	rows, err := s.dbPool.Query(context.Background(),
		`SELECT p.id, p.pub_time, p.link, p.title
		FROM posts p
		ORDER BY p.pub_time DESC
		LIMIT $1 OFFSET $2;`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.PostShort
	for rows.Next() {
		var p storage.PostShort
		err := rows.Scan(&p.ID, &p.PubTime, &p.Link, &p.Title)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Получает самые свежие новости из хранилища с учетом поиска
func (s *Storage) FindPosts(search string, offset, limit int) ([]storage.PostShort, error) {
	rows, err := s.dbPool.Query(context.Background(),
		`SELECT p.id, p.pub_time, p.link, p.title
		FROM posts p
		WHERE p.title LIKE concat('%', $1, '%')
		ORDER BY p.pub_time DESC
		LIMIT $2 OFFSET $3;`,
		search,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.PostShort
	for rows.Next() {
		var p storage.PostShort
		err := rows.Scan(&p.ID, &p.PubTime, &p.Link, &p.Title)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Получает общее количество новостей из хранилища
func (s *Storage) GetCount() (int, error) {
	var count int
	row := s.dbPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM posts;`)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Получает количество новостей из хранилища с учетом поиска
func (s *Storage) FindCount(search string) (int, error) {
	var count int
	row := s.dbPool.QueryRow(context.Background(),
		`SELECT COUNT(*)
		FROM posts p
		WHERE p.title LIKE concat('%', $1, '%');`,
		search,
	)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
