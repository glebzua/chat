package middlewares

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/http/controllers"
	"errors"
	"net/http"
)

type Userable interface {
	GetUserId() int64
	GetrecipientId() int64
}

func HaveAccessMiddleware[domainType Userable](key string, _ domainType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := ctx.Value(controllers.UserKey).(domain.User)
			obj := ctx.Value(key).(domainType)
			if obj.GetUserId() != user.Id && obj.GetrecipientId() != user.Id {
				err := errors.New("you have no access to this object")
				controllers.Forbidden(w, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func IsOwnerMiddleware[domainType Userable](key string, _ domainType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := ctx.Value(controllers.UserKey).(domain.User)
			obj := ctx.Value(key).(domainType)
			if obj.GetUserId() != user.Id {
				err := errors.New("you have no access to this object")
				controllers.Forbidden(w, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
