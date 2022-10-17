package http

import (
	"chatprjkt/config"
	"chatprjkt/config/container"
	"chatprjkt/internal/app"
	"chatprjkt/internal/infra/http/controllers"
	"chatprjkt/internal/infra/http/middlewares"

	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Router(cont container.Container) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(authRouter chi.Router) {
				authRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController, cont.UserService)
				ContactsRouter(apiRouter, cont.ContactsController, cont.UserService)
				MessagesRouter(apiRouter, cont.MessagesController, cont.UserService)

				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/change-pwd",
			ac.ChangePassword(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController, us app.UserService) {
	uom := middlewares.PathObject("id", controllers.PathUserKey, us)
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/me",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
		apiRouter.With(uom).Get(
			"/{id}",
			uc.FindOne(),
		)
		apiRouter.Get(
			"/",
			uc.FindAll(),
		)
	})
}

func ContactsRouter(r chi.Router, cc controllers.ContactsController, us app.UserService) {
	uom := middlewares.PathObject("id", controllers.PathUserKey, us)
	r.Route("/contacts", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/my",
			cc.FindAllMy(),
		)
		apiRouter.With(uom).Delete(
			"/{id}",
			cc.Delete(),
		)
		apiRouter.With(uom).Get(
			"/{id}",
			cc.FindOne(),
		)
		apiRouter.Post(
			"/",
			cc.Create(),
		)
		apiRouter.Get(
			"/",
			cc.FindAll(),
		)
	})
}
func MessagesRouter(r chi.Router, cc controllers.MessagesController, us app.UserService) {
	uom := middlewares.PathObject("id", controllers.PathUserKey, us)
	r.Route("/messages", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/my",
			cc.FindAllMy(),
		)
		apiRouter.With(uom).Delete(
			"/{id}",
			cc.Delete(),
		)
		apiRouter.With(uom).Get(
			"/{id}",
			cc.FindOne(),
		)
		apiRouter.Post(
			"/",
			cc.Create(),
		)
		apiRouter.Get(
			"/",
			cc.FindAll(),
		)
	})
}
func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("OK")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
