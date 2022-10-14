package container

import (
	"chatprjkt/config"
	"chatprjkt/internal/app"
	"chatprjkt/internal/infra/database"
	"chatprjkt/internal/infra/http/controllers"
	"chatprjkt/internal/infra/http/middlewares"
	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"net/http"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw func(http.Handler) http.Handler
}

type Services struct {
	app.UserService
	app.AuthService
	app.ContactsService
}

type Controllers struct {
	controllers.UserController
	controllers.AuthController
	controllers.ContactsController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	contactsRepository := database.NewContactsRepository(sess)
	sessionRepository := database.NewSessRepository(sess)

	userService := app.NewUserService(userRepository)
	contactsService := app.NewContactsService(contactsRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	contactsController := controllers.NewContactsController(contactsService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			userService,
			authService,
			contactsService,
		},
		Controllers: Controllers{
			userController,
			authController,
			contactsController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
