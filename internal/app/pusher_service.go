package app

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/database"
	"chatprjkt/internal/infra/resources"
)

type PusherService interface {
	NewUser(dto resources.UserDto)
	NewMessage(dto resources.MessageDto)
	UnreadMessages(dto domain.Messages)
}

type pusherService struct {
	pusherRepo database.PusherRepository
}

func (p pusherService) NewUser(dto resources.UserDto) {
	p.pusherRepo.NewUser(dto)
}
func (p pusherService) NewMessage(dto resources.MessageDto) {
	p.pusherRepo.NewMessage(dto)
}
func (p pusherService) UnreadMessages(dto domain.Messages) {
	p.pusherRepo.UnreadMessages(dto)
}

func NewPusherService(pr database.PusherRepository) PusherService {
	return pusherService{
		pusherRepo: pr,
	}
}
