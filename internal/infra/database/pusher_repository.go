package database

import (
	"chatprjkt/internal/infra/resources"
	"fmt"
	"github.com/pusher/pusher-http-go/v5"
)

//const pusherTableName = "pushers"

type PusherRepository interface {
	Save(dto resources.UserDto)
}

type pusherRepository struct {
	pusherConfig pusher.Client
}

func NewPusherRepository(cf pusher.Client) PusherRepository {
	return pusherRepository{
		pusherConfig: cf,
	}
}
func (p pusherRepository) Save(dto resources.UserDto) {

	client := p.pusherConfig

	data := map[string]string{"message from pusherRepository": fmt.Sprintf("add to contacts new user %s ", dto.Name)}
	err := client.Trigger("my-channel", "my-event", data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
