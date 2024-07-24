package api

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"sf-news/pkg/storage"

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
	a.router.HandleFunc("/news/{id}", a.handleGetPost).Methods(http.MethodGet)
	a.router.HandleFunc("/news", a.handleGetPosts).Methods(http.MethodGet)
}

// Обработчик для запроса новости по id в подробном виде
func (a *API) handleGetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, err := a.storage.GetPost(id)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// Обработчик для запроса новостей с возможностью фильтрации и пагенации
func (a *API) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	// Проверка номера страницы, по умолчанию - первая страница
	pageParam := r.URL.Query().Get("page")
	var page int
	var err error
	if pageParam == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if page <= 0 {
			http.Error(w, "positive page number required", http.StatusBadRequest)
			return
		}
	}

	// Проверка наличия строки поиска
	search := r.URL.Query().Get("search")
	var postsCount int
	var posts []storage.PostShort
	if search == "" {
		// Если строки поиска нет - то количесто и пагенацию ведем по всем новостям
		postsCount, err = a.storage.GetCount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts, err = a.storage.GetPosts((page-1)*ELEM_PER_PAGE, ELEM_PER_PAGE)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Если строка поиска есть - то используем её при подсчете количества новостей
		// и для работы пагенации
		postsCount, err = a.storage.FindCount(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts, err = a.storage.FindPosts(search, (page-1)*ELEM_PER_PAGE, ELEM_PER_PAGE)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возврат результата
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PostsCollection{
		Data:      posts,
		ElemCount: ELEM_PER_PAGE,
		PageCount: int(math.Ceil(float64(postsCount) / float64(ELEM_PER_PAGE))),
		PageNo:    page,
	})
}
