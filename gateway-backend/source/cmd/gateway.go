package main

import (
	"log"
	"net/http"
	"os"

	"sf-gateway/pkg/api"
	"sf-gateway/pkg/mdl"
	"sf-gateway/pkg/output"
)

func main() {
	// Создание логгера
	out := output.Make(os.Stdout, os.Stderr)

	// Создание роутера
	api := api.New()

	// Добавление middleware
	var router http.Handler = api.Router()
	router = mdl.WrapWithLogger(router, out)
	router = mdl.WrapWithId(router)
	router = mdl.WrapWithPingEcho(router)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err)
	}
}
