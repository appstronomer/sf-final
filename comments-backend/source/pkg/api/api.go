package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"sf-comments/pkg/storage"

	"github.com/gorilla/mux"
)

const ELEM_PER_PAGE = 10

type API struct {
	storage storage.StorageIface
	router  *mux.Router
}

// Конструктор объекта программного интерфейса, где storage - объект хранилища
func New(storage storage.StorageIface) *API {
	a := &API{
		storage: storage,
		router:  mux.NewRouter(),
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
	a.router.HandleFunc("/comments/post/{post_id}", a.handleGetComments).Methods(http.MethodGet)
	a.router.HandleFunc("/comments", a.handleAddComment).Methods(http.MethodPost)
}

// Размещает новый комментарий в хранилище
func (a *API) handleAddComment(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса
	var c storage.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление текущего времени
	c.PubTime = time.Now().Unix()

	// Запрос на добавление комментария в хранилище
	err = a.storage.PushComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат результата
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Получает комментарии для публикации
func (a *API) handleGetComments(w http.ResponseWriter, r *http.Request) {
	// Обязательный параметр пути post_id
	postId, err := strconv.Atoi(mux.Vars(r)["post_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if postId <= 0 {
		http.Error(w, "post_id path param should be a positive int", http.StatusBadRequest)
		return
	}

	// Опциональный query-параметр parent_id
	parentIdParam := r.URL.Query().Get("parent_id")
	var parentId int
	if parentIdParam == "" {
		parentId = 0
	} else {
		parentId, err = strconv.Atoi(parentIdParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if parentId <= 0 {
			http.Error(w, "positive parent_id required", http.StatusBadRequest)
			return
		}
	}

	// Опциональный query-параметр last_id
	lastIdParam := r.URL.Query().Get("last_id")
	var lastId int
	if lastIdParam == "" {
		lastId = 0
	} else {
		lastId, err = strconv.Atoi(lastIdParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if lastId <= 0 {
			http.Error(w, "positive last_id required", http.StatusBadRequest)
			return
		}
	}

	// Запрос в хранилище
	comments, err := a.storage.GetComments(postId, parentId, lastId, ELEM_PER_PAGE)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат результата
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
