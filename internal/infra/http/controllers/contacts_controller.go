package controllers

import (
	"chatprjkt/internal/infra/requests"
	"chatprjkt/internal/infra/resources"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"

	"chatprjkt/internal/app"
	"chatprjkt/internal/domain"
)

type ContactsController struct {
	contactsService app.ContactsService
}

func NewContactsController(s app.ContactsService) ContactsController {
	return ContactsController{
		contactsService: s,
	}
}

func (c ContactsController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.Contact)
		sess := r.Context().Value(SessKey).(domain.Session)

		err := (c.contactsService).Delete(u.Id, sess)
		if err != nil {
			log.Printf("ContactsController: %s", err)
			InternalServerError(w, err)
			return
		}

		ok(w)
	}
}
func (c ContactsController) FindAllMy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		cnt, err := (c.contactsService).FindAllForId(user.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		var cntDto resources.ContactDto
		success(w, cntDto.DomainToDtoCollection(cnt))

	}
}

//func (c ContactsController) FindAll() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cnt, err := (c.contactsService).FindAll()
//		if err != nil {
//			log.Print(err)
//			InternalServerError(w, err)
//			return
//		}
//		var cntDto resources.ContactDto
//		success(w, cntDto.DomainToDtoCollection(cnt))
//
//	}
//}

//func (c ContactsController) FindAll() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
//		if err != nil {
//			log.Println("contactsController FindAll error: ", err)
//			BadRequest(w, err)
//			return
//		}
//
//		cnt, err := (c.contactsService).FindAll(20, uint(page))
//		if err != nil {
//			log.Print(err)
//			InternalServerError(w, err)
//			return
//		}
//		var cntDto resources.ContactDto
//		success(w, cntDto.DomainToDtoCollection(cnt))
//
//	}
//}

func (c ContactsController) FindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contact := r.Context().Value(PathUserKey).(domain.User)
		cnt, err := (c.contactsService).FindById(contact.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		success(w, resources.ContactInfoDto{}.DomainToDto(cnt))
	}
}

func (c ContactsController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		cnt, err := requests.Bind(r, requests.ContactRequest{}, domain.Contact{})
		if err != nil {
			log.Printf("contactsController Create: %s", err)
			BadRequest(w, err)
			return
		}
		log.Println(cnt)
		cnt.UserId = user.Id
		log.Println(cnt.Nickname)

		alreadyInList, err := (c.contactsService).FindAllForId(cnt.UserId)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		for i := 0; i < len(alreadyInList.Items); i++ {
			if alreadyInList.Items[i].ContactId == cnt.ContactId {
				log.Print("already in list of contacts or waiting for confirm request")
				err = fmt.Errorf("already in list of contacts or waiting for confirm request")

				BadRequest(w, err)
				return
			}

		}
		activated, err := (c.contactsService).FindAllForId(cnt.ContactId)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		for i := 0; i < len(activated.Items); i++ {
			if activated.Items[i].ContactId == cnt.UserId {

				//c.contactsService.Activation(activated.Items[i].Id, cnt.Id)
				cnt.Activated = true

				chatId, err := c.GenerateChatIdHash(fmt.Sprintf("%v%v%v", cnt.UserId, time.Now(), activated.Items[i].UserId))
				if err != nil {
					log.Print(err)
					InternalServerError(w, err)
					return
				}
				err = c.contactsService.Activation(activated.Items[i].Id, chatId)
				if err != nil {
					log.Print(err)
					InternalServerError(w, err)
					return
				}
				cnt.ChatId = chatId
			}

		}
		cnt, err = c.contactsService.Save(cnt)
		if err != nil {
			log.Printf("ContactsController Create: %s", err)
			InternalServerError(w, err)
			return
		}
		var cntDto resources.ContactDto
		created(w, cntDto.DomainToDto(cnt))
	}
}
func (c ContactsController) GenerateChatIdHash(chatId string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(chatId), bcrypt.DefaultCost)
	return string(bytes), err
}
