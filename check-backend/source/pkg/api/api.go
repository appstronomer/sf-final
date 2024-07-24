package api

import (
	"encoding/json"
	"net/http"
	"sf-check/pkg/checker"

	"github.com/gorilla/mux"
)

const ELEM_PER_PAGE = 10

type API struct {
	router *mux.Router
}

// Конструктор объекта программного интерфейса, где storage - объект хранилища
func New() *API {
	a := &API{
		router: mux.NewRouter(),
	}
	a.endpoints()
	return a
}

// Возвращает внутренний роутер
func (a *API) Router() *mux.Router {
	return a.router
}

// Регистрирует все обработчики
func (a *API) endpoints() {
	a.router.HandleFunc("/check", a.handleCheckComment).Methods(http.MethodPost)
}

// Размещает новый комментарий в хранилище
func (a *API) handleCheckComment(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса
	var c checker.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if checker.CheckIfIncorrect(c) {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
