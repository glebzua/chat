package middlewares

import (
	"chatprjkt/internal/infra/http/controllers"
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/upper/db/v4"
)

type Findable interface {
	Find(int64) (interface{}, error)
}

func PathObject(pathKey string, ctxKey string, service Findable) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseUint(chi.URLParam(r, pathKey), 10, 64)
			if err != nil {
				err = errors.New("invalid user id parameter(only non-negative integers)")
				log.Print(err)
				controllers.BadRequest(w, err)
				return
			}

			obj, err := service.Find(int64(id))
			if err != nil {
				if err == db.ErrNoMoreRows {
					log.Print(err)
					controllers.NotFound(w, err)
					return
				}
				log.Print(err)
				controllers.InternalServerError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey, obj)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
