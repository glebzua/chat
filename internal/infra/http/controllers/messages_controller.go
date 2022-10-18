package controllers

import (
	"chatprjkt/internal/app"
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/requests"
	"chatprjkt/internal/infra/resources"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type MessagesController struct {
	messagesService app.MessagesService
}

func NewMessagesController(s app.MessagesService) MessagesController {
	return MessagesController{
		messagesService: s,
	}
}

func (c MessagesController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.Contact)
		sess := r.Context().Value(SessKey).(domain.Session)

		err := (c.messagesService).Delete(u.Id, sess)
		if err != nil {
			log.Printf("MessagesController: %s", err)
			InternalServerError(w, err)
			return
		}

		ok(w)
	}
}
func (c MessagesController) FindAllMy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		msg, err := (c.messagesService).FindAllForId(user.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		var msgDto resources.MessageDto
		success(w, msgDto.DomainToDtoCollection(msg))

	}
}
func (c MessagesController) FindAllMyChats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		msg, err := (c.messagesService).FindAllForId(user.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		var chats domain.Chats
		var chat domain.Chat
		keys := make(map[string]bool)
		var uniq []any

		for _, entry := range msg.Items {
			if _, value := keys[entry.ChatId]; !value {
				keys[entry.ChatId] = true
				uniq = append(uniq, entry.ChatId, entry.RecipientId)
				log.Println("string(entry.RecipientId)", entry.RecipientId)
				log.Println("entry.RecipientId)", entry.RecipientId)
				log.Println("entry.ChatId)", entry.ChatId)
				chat.ChatId = entry.ChatId
				chat.RecipientId = entry.RecipientId
				chats.Items = append(chats.Items, chat)
			}
		}
		log.Println("uniq=", uniq)
		var ChatDto resources.ChatDto
		success(w, ChatDto.DomainToDtoCollection(chats))

	}
}

func (c MessagesController) FindAllMessagesInChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := r.URL.Query().Get("chatId")
		if chatId == "" {
			log.Println("messagesController FindAll error: ", chatId)
			BadRequest(w, fmt.Errorf(" messagesController: %v", chatId))
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		msg, err := (c.messagesService).FindAllMessagesInChat(user.Id, chatId)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}

		var MessageDto resources.MessageDto
		success(w, MessageDto.DomainToDtoCollection(msg))

	}
}

func (c MessagesController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			log.Println("MessagesController FindAll error: ", err)
			BadRequest(w, err)
			return
		}

		msg, err := (c.messagesService).FindAll(20, uint(page))
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		var msgDto resources.MessageDto
		success(w, msgDto.DomainToDtoCollection(msg))

	}
}
func (c MessagesController) FindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contact := r.Context().Value(PathUserKey).(domain.User)
		cnt, err := (c.messagesService).FindById(contact.Id)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}
		success(w, resources.MessageDto{}.DomainToDto(cnt))
	}
}

func (c MessagesController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		msg, err := requests.Bind(r, requests.MessagesRequest{}, domain.Message{})
		if err != nil {
			log.Printf("MessagesController Create: %s", err)
			BadRequest(w, err)
			return
		}

		msg.SenderId = user.Id

		msg, err = c.messagesService.Save(msg)
		if err != nil {
			log.Printf("MessagesController Create: %s", err)
			InternalServerError(w, err)
			return
		}
		var msgDto resources.MessageDto
		created(w, msgDto.DomainToDto(msg))
	}
}
