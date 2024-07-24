package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sf-check/pkg/mdl"
	"sf-check/pkg/model"
	"strconv"
	"sync"

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
	a.router.HandleFunc("/news/latest", a.handleNewsLatest).Methods(http.MethodGet)
	a.router.HandleFunc("/news/{id}", a.handleNewsById).Methods(http.MethodGet)
	a.router.HandleFunc("/comments/post/{post_id}", a.handleGetComments).Methods(http.MethodGet)
	a.router.HandleFunc("/comments", a.handlePostComment).Methods(http.MethodPost)
}

// Запрос на размещение комментария
func (a *API) handlePostComment(w http.ResponseWriter, r *http.Request) {
	query := url.Values{}

	// Добавление сквозного идентификатора запроса
	requestId, ok := r.Context().Value(mdl.MdlKey("request_id")).(string)
	if ok {
		query.Add("request_id", requestId)
	}

	// Валидация тела запроса
	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commentBytes, err := json.Marshal(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка комментария на стоп-слова
	res, err := http.Post(
		fmt.Sprintf("http://check-backend/check?%s", query.Encode()),
		"application/json",
		bytes.NewReader(commentBytes),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode == http.StatusBadRequest {
		http.Error(w, "incorrect comment", http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	// Отправка комментария в хранилище
	res, err = http.Post(
		fmt.Sprintf("http://comments-backend/comments?%s", query.Encode()),
		"application/json",
		bytes.NewReader(commentBytes),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Проброс ошибки клиенту
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, res.Body)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Запрос комментария из древа комментариев
func (a *API) handleGetComments(w http.ResponseWriter, r *http.Request) {
	query := url.Values{}

	// Добавление сквозного идентификатора запроса
	requestId, ok := r.Context().Value(mdl.MdlKey("request_id")).(string)
	if ok {
		query.Add("request_id", requestId)
	}

	// Получение обязательного ID новости
	newsId, err := strconv.Atoi(mux.Vars(r)["post_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление опционального параметра parent_id
	paramName := "parent_id"
	err = addQueryIntParam(r, query, paramName, paramName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление опционального параметра last_id
	paramName = "last_id"
	err = addQueryIntParam(r, query, paramName, paramName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создание запроса
	res, err := http.Get(fmt.Sprintf("http://comments-backend/comments/post/%d?%s", newsId, query.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Проброс ошибки клиенту
	newsBytes, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		w.Write(newsBytes)
		return
	}

	// Проброс успеха клиенту
	w.Header().Set("Content-Type", "application/json")
	w.Write(newsBytes)
}

// Отдает конкретную новость с первой страницей комментариев
func (a *API) handleNewsById(w http.ResponseWriter, r *http.Request) {
	queryNewsItem := url.Values{}
	queryComments := url.Values{}

	// Добавление сквозного идентификатора запроса
	requestId, ok := r.Context().Value(mdl.MdlKey("request_id")).(string)
	if ok {
		queryNewsItem.Add("request_id", requestId)
		queryComments.Add("request_id", requestId)
	}

	// Получение обязательного ID новости
	newsId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Элементы синхронизации
	ch := make(chan interface{}, 2)
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Запуск двух запросов параллельно
	go getNewsItem(&wg, ch, newsId, queryNewsItem)
	go getComments(&wg, ch, newsId, queryComments)
	wg.Wait()

	// Интерпретация результатов
	var result model.NewsComplex
	for item := range ch {
		switch itemConc := item.(type) {
		case model.NewsFullDetailed:
			result.Data = itemConc
		case []model.Comment:
			result.Comments = itemConc
		case error:
			http.Error(w, itemConc.Error(), http.StatusInternalServerError)
			return
		default:
			http.Error(w, "unexected result", http.StatusInternalServerError)
			return
		}
	}

	// Отправка результата потребителю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Отдает список новостей с учетом пагенацией и опциональной фильтрации
func (a *API) handleNewsLatest(w http.ResponseWriter, r *http.Request) {
	query := url.Values{}

	// Добавление сквозного идентификатора запроса
	requestId, ok := r.Context().Value(mdl.MdlKey("request_id")).(string)
	if ok {
		query.Add("request_id", requestId)
	}

	// Добавление опционального номера страницы
	paramName := "page"
	err := addQueryIntParam(r, query, paramName, paramName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: удалить после тестов
	// pageParam := r.URL.Query().Get("page")
	// var page int
	// var err error
	// if pageParam != "" {
	// 	page, err = strconv.Atoi(pageParam)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	if page <= 0 {
	// 		http.Error(w, "positive page number required", http.StatusBadRequest)
	// 		return
	// 	}
	// 	query.Add("page", strconv.Itoa(page))
	// }

	// Добалвение опциональной строки поиска
	search := r.URL.Query().Get("search")
	if search != "" {
		query.Add("search", search)
	}

	// Создание запроса
	res, err := http.Get(fmt.Sprintf("http://news-backend/news?%s", query.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Проброс ошибки клиенту
	newsBytes, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		w.Write(newsBytes)
		return
	}

	// Валидация ответа сервера новостей
	var newsCollection model.NewsCollection
	err = json.Unmarshal(newsBytes, &newsCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка результата потребителю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsCollection)
}

func getNewsItem(wg *sync.WaitGroup, ch chan<- interface{}, newsId int, queryNewsItem url.Values) {
	defer wg.Done()
	res, err := http.Get(fmt.Sprintf("http://news-backend/news/%d?%s", newsId, queryNewsItem.Encode()))
	if err != nil {
		ch <- err
		return
	}
	defer res.Body.Close()
	ch <- res
}

func getComments(wg *sync.WaitGroup, ch chan<- interface{}, newsId int, queryComments url.Values) {
	defer wg.Done()
	res, err := http.Get(fmt.Sprintf("http://comments-backend/comments/post/%d?%s", newsId, queryComments.Encode()))
	if err != nil {
		ch <- err
		return
	}
	defer res.Body.Close()
	ch <- res
}

func addQueryIntParam(r *http.Request, query url.Values, upParamName, downParamName string) error {
	valParam := r.URL.Query().Get(upParamName)
	var val int
	var err error
	if valParam != "" {
		val, err = strconv.Atoi(valParam)
		if err != nil {
			return err
		}
		if val <= 0 {
			return err
		}
		query.Add(downParamName, strconv.Itoa(val))
	}
	return nil
}
