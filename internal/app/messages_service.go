package app

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/database"
	"fmt"
	"log"
)

type MessagesService interface {
	FindAll(pageSize, page uint) (domain.Messages, error)
	FindAllForId(id int64) (domain.Messages, error)
	FindAllMessagesInChat(userId int64, chatId string) (domain.Messages, error)
	Save(item domain.Message) (domain.Message, error)
	Find(id int64) (interface{}, error)
	FindById(id int64) (domain.Message, error)

	Delete(id int64, sess domain.Session) error
}

type messagesService struct {
	messagesRepo    database.MessagesRepository
	contactsService ContactsService
}

func NewMessagesService(mr database.MessagesRepository, cs ContactsService) MessagesService {
	return messagesService{
		messagesRepo:    mr,
		contactsService: cs,
	}
}
func (s messagesService) Save(item domain.Message) (domain.Message, error) {
	var result domain.Message
	contacts, err := s.contactsService.FindAllForId(item.SenderId)
	if len(contacts.Items) == 0 {
		err = fmt.Errorf("user have no contacts")
		return domain.Message{}, err
	}
	for _, contact := range contacts.Items {
		if (contact.ContactId == item.RecipientId) && (contact.ChatId == item.ChatId) {
			result, err := s.messagesRepo.Save(item)
			if err != nil {
				return domain.Message{}, fmt.Errorf("messagesService Save: %w", err)
			}
			return result, err
		}
	}
	err = fmt.Errorf("recipient or ChatId incorrect")
	return result, err
}

func (s messagesService) FindAll(pageSize, page uint) (domain.Messages, error) {
	contacts, err := s.messagesRepo.FindAll(pageSize, page)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return domain.Messages{}, err
	}

	return contacts, nil
}
func (s messagesService) FindAllForId(id int64) (domain.Messages, error) {
	contacts, err := s.contactsService.FindAllForId(id)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return domain.Messages{}, err
	}
	var messages domain.Messages
	for _, chat := range contacts.Items {
		message, err := s.messagesRepo.FindAllForId(chat.ChatId)
		if err != nil {
			log.Printf("messagesService: %s", err)
			return domain.Messages{}, err
		}
		messages.Items = append(messages.Items, message.Items...)
	}
	return messages, nil
}
func (s messagesService) FindAllMessagesInChat(id int64, chatId string) (domain.Messages, error) {
	messages, err := s.messagesRepo.FindAllMessagesInChat(id, chatId)
	if err != nil {
		log.Printf("messagesService: %s", err)
		return domain.Messages{}, err
	}
	//var messages domain.Messages
	//for _, chat := range contacts.Items {
	//	message, err := s.messagesRepo.FindAllForId(chat.ChatId)
	//	if err != nil {
	//		log.Printf("messagesService: %s", err)
	//		return domain.Messages{}, err
	//	}
	//	messages.Items = append(messages.Items, message.Items...)
	//}
	return messages, nil
}

func (s messagesService) FindById(id int64) (domain.Message, error) {
	contact, err := s.messagesRepo.FindById(id)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return domain.Message{}, err
	}

	return contact, err
}
func (s messagesService) Find(id int64) (interface{}, error) {
	return s.messagesRepo.FindById(id)
}
func (s messagesService) Delete(id int64, sess domain.Session) error {
	err := s.messagesRepo.Delete(id)
	if err != nil {
		log.Printf("messageService: %s", err)
		return err
	}

	return nil
}
