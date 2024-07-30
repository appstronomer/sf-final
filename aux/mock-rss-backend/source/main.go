package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Путь к дирректории со статическими xml-файлами должен быть передан
	// в качестве аргумента
	if len(os.Args) != 2 {
		log.Fatal("please, provide path to static web directory as the only argument")
	}

	// Запускается сервер статических файлов
	fs := http.FileServer(http.Dir(os.Args[1]))
	http.Handle("/", fs)
	log.Println("Listening on :80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
