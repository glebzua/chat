package controllers

import (
	"chatprjkt/internal/app"
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/requests"
	"chatprjkt/internal/infra/resources"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(s app.UserService) UserController {
	return UserController{
		userService: s,
	}
}

func (c UserController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.UserRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
		}

		user, err = (c.userService).Save(user)
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}
		//
		var userDto resources.UserDto
		created(w, userDto.DomainToDto(user))
	}
}

func (c UserController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.UserRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}

		u := r.Context().Value(UserKey).(domain.User)
		user, err = c.userService.Update(u, user)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		var userDto resources.UserDto
		success(w, userDto.DomainToDto(user))
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		sess := r.Context().Value(SessKey).(domain.Session)

		err := (c.userService).Delete(u.Id, sess)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		ok(w)
	}
}
func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		usr, err := (c.userService).FindById(user.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		success(w, resources.UserDto{}.DomainToDto(usr))
	}
}

func (c UserController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			log.Println("Client FindAll error: ", err)
			BadRequest(w, err)
			return
		}
		usr, err := (c.userService).FindAll(20, uint(page))
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		success(w, resources.UserInfoDto{}.DomainToDtoCollection(usr))

	}
}
func (c UserController) FindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(PathUserKey).(domain.User)

		usr, err := (c.userService).FindById(user.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		success(w, resources.UserInfoDto{}.DomainToDto(usr))
	}
}
