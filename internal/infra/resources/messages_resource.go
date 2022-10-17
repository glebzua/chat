package resources

import (
	"chatprjkt/internal/domain"
)

type MessageDto struct {
	Id          int64  `json:"id"`
	ChatId      string `json:"chatId"`
	SenderId    int64  `json:"senderId"`
	RecipientId int64  `json:"recipientId"`
	Message     string `json:"message"`
	Sended      bool   `json:"sended"`
	Received    bool   `json:"received"`
}

type MessagesDto struct {
	message []MessageDto
}

func (d MessageDto) DomainToDto(message domain.Message) MessageDto {
	return MessageDto{
		Id:          message.Id,
		ChatId:      message.ChatId,
		SenderId:    message.SenderId,
		RecipientId: message.RecipientId,
		Message:     message.Message,
		Sended:      message.Sended,
		Received:    message.Received,
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
