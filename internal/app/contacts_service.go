package app

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/database"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type ContactsService interface {
	//FindAll() (domain.Contacts, error)
	//FindAll(pageSize, page uint) (domain.Contacts, error)
	FindAllForId(id int64) (domain.Contacts, error)
	Save(item domain.Contact) (domain.Contact, error)
	Find(id int64) (interface{}, error)
	FindById(id int64) (domain.Contact, error)

	Delete(id int64, sess domain.Session) error
	//Activation(id int64, id2 int64) error
	Activation(id int64, chatid string) error
}

type contactsService struct {
	contactsRepo database.ContactsRepository
}

func NewContactsService(cr database.ContactsRepository) ContactsService {
	return contactsService{
		contactsRepo: cr,
	}
}
func (s contactsService) Save(item domain.Contact) (domain.Contact, error) {
	result, err := s.contactsRepo.Save(item)
	if err != nil {
		return domain.Contact{}, fmt.Errorf("service Save: %w", err)
	}
	return result, err
}

//func (s contactsService) FindAll() (domain.Contacts, error) {
//	contacts, err := s.contactsRepo.FindAll()
//	if err != nil {
//		log.Printf("contactsService: %s", err)
//		return domain.Contacts{}, err
//	}
//
//	return contacts, nil
//}

//	func (s contactsService) FindAll(pageSize, page uint) (domain.Contacts, error) {
//		contacts, err := s.contactsRepo.FindAll(pageSize, page)
//		if err != nil {
//			log.Printf("contactsService: %s", err)
//			return domain.Contacts{}, err
//		}
//
//		return contacts, nil
//	}
func (s contactsService) FindAllForId(id int64) (domain.Contacts, error) {
	contacts, err := s.contactsRepo.FindAllForId(id)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return domain.Contacts{}, err
	}

	return contacts, nil
}
func (s contactsService) Activation(id int64, chatid string) error {
	err := s.contactsRepo.Activation(id, chatid)
	//err := s.contactsRepo.Activation(id, id2)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return err
	}

	return err
}
func (s contactsService) FindById(id int64) (domain.Contact, error) {
	contact, err := s.contactsRepo.FindById(id)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return domain.Contact{}, err
	}

	return contact, err
}
func (s contactsService) Find(id int64) (interface{}, error) {
	return s.contactsRepo.FindById(id)
}
func (s contactsService) Delete(id int64, sess domain.Session) error {
	err := s.contactsRepo.Delete(id)
	if err != nil {
		log.Printf("contactsService: %s", err)
		return err
	}

	return nil
}

func (s contactsService) GenerateChatIdHash(chatId string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(chatId), bcrypt.DefaultCost)
	return string(bytes), err
}
