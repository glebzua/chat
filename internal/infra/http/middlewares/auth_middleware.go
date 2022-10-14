package middlewares

import (
	"chatprjkt/internal/app"
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/http/controllers"
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

func AuthMiddleware(ja *jwtauth.JWTAuth, as app.AuthService, us app.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := jwtauth.VerifyRequest(ja, r, jwtauth.TokenFromHeader)

			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				controllers.Unauthorized(w, err)
				return
			}

			claims := token.PrivateClaims()
			uId := int64(claims["user_id"].(float64))
			uUuid, err := uuid.Parse(claims["uuid"].(string))
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			auth := domain.Session{
				UserId: uId,
				UUID:   uUuid,
			}
			err = as.Check(auth)
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			user, err := us.FindById(uId)
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			ctx = context.WithValue(ctx, controllers.UserKey, user)
			ctx = context.WithValue(ctx, controllers.SessKey, auth)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
