package resources

import (
	"chatprjkt/internal/domain"
	"time"
)

type ContactDto struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	ContactId   int64     `json:"contactId"`
	Activated   bool      `json:"activated"`
	ChatId      string    `json:"chatId"`
	Nickname    string    `json:"nickname"`
	CreatedDate time.Time `json:"createdDate"`
}
type ContactInfoDto struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	ContactId   int64     `json:"contactId"`
	Activated   bool      `json:"activated"`
	ChatId      string    `json:"chatId"`
	Nickname    string    `json:"nickname"`
	CreatedDate time.Time `json:"createdDate"`
}

type ContactsDto struct {
	contact []ContactDto `json:"contact"`
}

func (d ContactDto) DomainToDto(contact domain.Contact) ContactDto {
	return ContactDto{
		Id:          contact.Id,
		UserId:      contact.UserId,
		ContactId:   contact.ContactId,
		Activated:   contact.Activated,
		ChatId:      contact.ChatId,
		Nickname:    contact.Nickname,
		CreatedDate: contact.CreatedDate,
	}
}

func (d ContactInfoDto) DomainToDto(contact domain.Contact) ContactInfoDto {
	return ContactInfoDto{
		Id:          contact.Id,
		UserId:      contact.UserId,
		ContactId:   contact.ContactId,
		Activated:   contact.Activated,
		ChatId:      contact.ChatId,
		Nickname:    contact.Nickname,
		CreatedDate: contact.CreatedDate,
	}
}

func (d ContactDto) DomainToDtoCollection(u domain.Contacts) []ContactDto {
	var result []ContactDto

	for _, c := range u.Items {
		var dto ContactDto
		contactDto := dto.DomainToDto(c)
		result = append(result, contactDto)
	}

	return result
}
