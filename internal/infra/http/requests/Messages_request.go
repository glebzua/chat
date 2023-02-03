package requests

import "chatprjkt/internal/domain"

type MessagesRequest struct {
	RecipientId int64  `json:"recipientId"`
	ChatId      string `json:"chatId"`
	Message     string `json:"message"`
}

func (r MessagesRequest) ToDomainModel() (interface{}, error) {
	return domain.Message{
		RecipientId: r.RecipientId,
		ChatId:      r.ChatId,
		Message:     r.Message,
	}, nil
}
