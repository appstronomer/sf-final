package parser

import (
	"net/http"
	"time"

	"sf-news/pkg/config"
	"sf-news/pkg/output"
	"sf-news/pkg/rss"
	"sf-news/pkg/storage"
	"sf-news/pkg/storage/postgres"
)

// ЧТЕНИЕ НОВОСТЕЙ
// Инициализация всех горутин, осуществляющих чтение новостей
func InitParser(out output.Output, cfg config.Config, db *postgres.Storage) {
	delay := time.Duration(cfg.RequestPeriod) * time.Minute
	for _, url := range cfg.RssUrls {
		go parseLoop(out, db, url, delay)
	}
}

// Цкл чтения новостей конкретной горутиной
func parseLoop(out output.Output, db *postgres.Storage, url string, delay time.Duration) {
	for {
		posts, err := parseUrl(url)
		if err != nil {
			out.Err(err)
			continue
		}
		err = db.PushPosts(posts)
		if err != nil {
			out.Err(err)
		}
		time.Sleep(delay)
	}
}

// Итерация получения и десереализации новостей
func parseUrl(url string) ([]storage.PostFull, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	posts, err := rss.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
