package app

import (
	"chatprjkt/internal/infra/database"
	"chatprjkt/internal/infra/resources"
)

type PusherService interface {
	Save(dto resources.UserDto)
}

type pusherService struct {
	pusherRepo database.PusherRepository
}

func (p pusherService) Save(dto resources.UserDto) {
	p.pusherRepo.Save(dto)
}

func NewPusherService(pr database.PusherRepository) PusherService {
	return pusherService{
		pusherRepo: pr,
	}
}
