package config

import (
	"encoding/json"
	"net/url"
	"os"
)

// КОНФИГУРАЦИЯ ПРИЛОЖЕНИЯ
// Контейнер под конфиг приложения
type Config struct {
	RssUrls       []string `json:"rss"`
	RequestPeriod uint32   `json:"request_period"`
}

// Чтение и валидация конфига
func ReadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return Config{}, err
	}
	for _, urlItem := range cfg.RssUrls {
		_, err := url.ParseRequestURI(urlItem)
		if err != nil {
			return Config{}, err
		}
	}
	return cfg, nil
}
