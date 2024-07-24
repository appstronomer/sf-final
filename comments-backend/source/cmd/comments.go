package main

import (
	"log"
	"net/http"
	"os"

	"sf-comments/pkg/api"
	"sf-comments/pkg/mdl"
	"sf-comments/pkg/output"
	"sf-comments/pkg/storage/postgres"
)

// Сообщение - мануал
const helpMessage string = `news DB
DB - postgres connection config string`

func main() {
	// Создание логгера
	out := output.Make(os.Stdout, os.Stderr)

	// Проверка аргументов запуска
	if len(os.Args) != 1 {
		log.Fatalf("ERROR: invalid arguments passed to application: %v\n\nSYNOPSIS\n%v", os.Args, helpMessage)
	}

	// Подключение к БД
	// Пример:"postgres://user:password@postgres:5432/sf"
	db, err := postgres.New(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Создание роутера
	api := api.New(db)

	// Добавление middleware
	var router http.Handler = api.Router()
	router = mdl.WrapWithId(router)
	router = mdl.WrapWithLogger(router, out)

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err)
	}
}
