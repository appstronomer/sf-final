package mdl

import (
	"fmt"
	"net/http"
	"sf-news/pkg/output"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func WrapWithLogger(handler http.Handler, out output.Output) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		handler.ServeHTTP(res, r)
		requestId, ok := r.Context().Value(MdlKey("request_id")).(string)
		if !ok {
			requestId = "undefined"
		}

		// Время запроса — фактическое время создания лога.
		// Уникальный ID запроса — содержится в объекте запроса.
		// IP-адрес, с которого был отправлен запрос — содержится в объекте запроса.
		// HTTP-код ответа — содержится в объекте ответа.
		out.Log(fmt.Sprintf("%s %s %s %d", time.Now().String(), requestId, r.RemoteAddr, res.statusCode))

	})
}
