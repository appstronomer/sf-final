package mdl

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type MdlKey string

func WrapWithId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.URL.Query().Get("request_id")
		if requestId == "" {
			requestId = fmt.Sprintf("cmb:%s", uuid.New().String())
		}
		ctx := context.WithValue(r.Context(), MdlKey("request_id"), requestId)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
