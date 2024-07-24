package main

import (
	"log"
	"net/http"
	"os"

	"sf-check/pkg/api"
	"sf-check/pkg/mdl"
	"sf-check/pkg/output"
)

func main() {
	// Создание логгера
	out := output.Make(os.Stdout, os.Stderr)

	// Создание роутера
	api := api.New()

	// Добавление middleware
	var router http.Handler = api.Router()
	router = mdl.WrapWithId(router)
	router = mdl.WrapWithLogger(router, out)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err)
	}
}
