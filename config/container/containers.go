package container

import (
	"chatprjkt/config"
	"chatprjkt/internal/app"
	"chatprjkt/internal/infra/database"
	"chatprjkt/internal/infra/filesystem"
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
	app.MessagesService
	app.PusherService
	app.ImageService
}

type Controllers struct {
	controllers.UserController
	controllers.AuthController
	controllers.ContactsController
	controllers.MessagesController
	controllers.PusherController
	controllers.ImageController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	contactsRepository := database.NewContactsRepository(sess)
	messagesRepository := database.NewMessagesRepository(sess)
	sessionRepository := database.NewSessRepository(sess)
	pusherRepository := database.NewPusherRepository(conf.Pusher)
	imageRepository := database.NewImageRepository(sess)

	pusherService := app.NewPusherService(pusherRepository)
	userService := app.NewUserService(userRepository, pusherService)
	contactsService := app.NewContactsService(contactsRepository)
	messagesService := app.NewMessagesService(messagesRepository, contactsService, pusherService)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)
	imageStorageService := filesystem.NewImageStorageService(conf.FileStorageLocation)
	imageService := app.NewImageService(imageRepository, imageStorageService, messagesService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	contactsController := controllers.NewContactsController(contactsService)
	messagesController := controllers.NewMessagesController(messagesService)
	pusherController := controllers.NewPusherController(pusherService)
	imageController := controllers.NewImageController(imageService)
	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			userService,
			authService,
			contactsService,
			messagesService,
			pusherService,
			imageService,
		},
		Controllers: Controllers{
			userController,
			authController,
			contactsController,
			messagesController,
			pusherController,
			imageController,
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
