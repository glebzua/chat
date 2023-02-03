package resources

import (
	"chatprjkt/internal/domain"
	"time"
)

type MessageDto struct {
	Id          int64      `json:"id"`
	ChatId      string     `json:"chatId"`
	SenderId    int64      `json:"senderId"`
	RecipientId int64      `json:"recipientId"`
	Message     string     `json:"message"`
	FileLoc     string     `json:"fileloc"`
	Sended      bool       `json:"sended"`
	Received    bool       `json:"received"`
	CreatedDate time.Time  `json:"createdDate"`
	UpdatedDate time.Time  `json:"updatedDate"`
	DeletedDate *time.Time `json:"deletedDate"`
}

type ChatDto struct {
	ChatId      string `json:"chatId"`
	RecipientId int64  `json:"recipientId"`
}

type MessagesDto struct {
	message []MessageDto
}
type ChatsDto struct {
	chat []ChatDto
}

func (d MessageDto) DomainToDto(message domain.Message) MessageDto {
	return MessageDto{
		Id:          message.Id,
		ChatId:      message.ChatId,
		SenderId:    message.SenderId,
		RecipientId: message.RecipientId,
		FileLoc:     message.FileLoc,
		Message:     message.Message,
		Sended:      message.Sended,
		Received:    message.Received,
		CreatedDate: message.CreatedDate,
		UpdatedDate: message.UpdatedDate,
		DeletedDate: message.DeletedDate,
	}
}

func (d ChatDto) DomainToDto(chat domain.Chat) ChatDto {
	return ChatDto{
		ChatId:      chat.ChatId,
		RecipientId: chat.RecipientId,
	}
}
func (d MessageDto) DomainToDtoCollection(u domain.Messages) []MessageDto {
	var result []MessageDto

	for _, c := range u.Items {
		var dto MessageDto
		messageDto := dto.DomainToDto(c)
		result = append(result, messageDto)
	}

	return result
}

func (d ChatDto) DomainToDtoCollection(u domain.Chats) []ChatDto {
	var result []ChatDto

	for _, c := range u.Items {
		var dto ChatDto
		chatDto := dto.DomainToDto(c)
		result = append(result, chatDto)
	}

	return result
}
