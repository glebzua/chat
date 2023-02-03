package controllers

import (
	"chatprjkt/internal/app"
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/http/resources"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ImageController struct {
	service app.ImageService
}

func NewImageController(s app.ImageService) ImageController {
	return ImageController{
		service: s,
	}
}

func (c ImageController) AddImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := r.URL.Query().Get("chatId")
		if chatId == "" {
			log.Println("ImageController AddImage  chatId error: ", chatId)
			BadRequest(w, fmt.Errorf(" imageController: %v", chatId))
			return
		}
		recipientId, err := strconv.ParseInt(r.URL.Query().Get("recipientId"), 10, 0)
		if err != nil {
			log.Println("ImageController AddImage  recipientId error: ", recipientId)
			BadRequest(w, fmt.Errorf(" imageController: %v", recipientId))
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		buff, err := io.ReadAll(r.Body)
		//buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ImageController AddImage user: %s", err)
			BadRequest(w, err)
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			err = errors.New("this file format is not allowed. upload JPEG or PNG")
			Forbidden(w, err)
			return
		}

		img := domain.Image{
			SenderId:    user.Id,
			RecipientId: recipientId,
			ChatId:      chatId,
			Name:        uuid.NewString() + "." + strings.TrimLeft(filetype, "image/"),
		}
		i, err := c.service.Save(img, buff)
		if err != nil {
			log.Printf("ImageController: %s", err)
			BadRequest(w, err)
			return
		}

		var imageDto resources.ImageDto
		fmt.Println("ImageController add image ", imageDto.DomainToDto(i))
		//createdImage(w, buff)
		created(w, imageDto.DomainToDto(i))
	}
}

func (c ImageController) DeleteImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		img := r.Context().Value(PathImgKey).(domain.Image)

		err := c.service.Delete(img.Id)
		if err != nil {
			log.Printf("ImageController: %s", err)
			InternalServerError(w, err)
			return
		}

		ok(w)
	}
}
func (c ImageController) FindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		img := r.Context().Value(PathImgKey).(domain.Image)

		i, err := c.service.FindByName(img.Name)
		if err != nil {
			log.Printf("imageController FindOne: %s", err)
			InternalServerError(w, err)
			return
		}
		var imageDto resources.ImageDto
		created(w, imageDto.DomainToDto(i))
	}
}
