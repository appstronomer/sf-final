package main

import (
	"log"
	"net/http"
	"os"

	"sf-news/pkg/api"
	"sf-news/pkg/config"
	"sf-news/pkg/mdl"
	"sf-news/pkg/output"
	"sf-news/pkg/parser"
	"sf-news/pkg/storage/postgres"
)

// Сообщение - мануал
const helpMessage string = `news CONFIG DB
CONFIG - path to application config.json file
DB     - postgres connection config string`

func main() {
	// Создание логгера
	out := output.Make(os.Stdout, os.Stderr)

	// Проверка аргументов запуска
	if len(os.Args) != 3 {
		log.Fatalf("ERROR: invalid arguments passed to application: %v\n\nSYNOPSIS\n%v", os.Args, helpMessage)
	}

	// Получение конфигурации
	cfg, err := config.ReadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к БД
	// Пример:"postgres://user:password@postgres:5432/sf"
	db, err := postgres.New(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// Создание роутера
	api := api.New(db)

	// Запуск парсеров
	parser.InitParser(out, cfg, db)

	// Добавление middleware
	var router http.Handler = api.Router()
	router = mdl.WrapWithId(router)
	router = mdl.WrapWithLogger(router, out)

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err)
	}
}
