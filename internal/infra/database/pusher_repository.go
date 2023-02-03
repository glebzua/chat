package database

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/http/resources"
	"fmt"
	"github.com/pusher/pusher-http-go/v5"
	"strconv"
)

type PusherRepository interface {
	NewUser(dto resources.UserDto)
	NewMessage(dto resources.MessageDto)
	UnreadMessages(dto domain.Messages)
}

type pusherRepository struct {
	pusherConfig pusher.Client
}

func NewPusherRepository(cf pusher.Client) PusherRepository {
	return pusherRepository{
		pusherConfig: cf,
	}
}
func (p pusherRepository) NewUser(dto resources.UserDto) {

	client := p.pusherConfig

	data := map[string]string{"new user registered": fmt.Sprintf(" %s", dto.Name)}
	err := client.Trigger("my-channel", "new-user", data)

	if err != nil {
		fmt.Println(err.Error())
	}
}
func (p pusherRepository) NewMessage(dto resources.MessageDto) {

	client := p.pusherConfig

	data := map[string]int64{"NewMessageFrom": dto.SenderId}
	err := client.Trigger(strconv.FormatInt(dto.RecipientId, 10), "NewMessage", data)

	if err != nil {
		fmt.Println("error in trigger", strconv.FormatInt(dto.RecipientId, 10))
		fmt.Println(err.Error())
	}
}

func (p pusherRepository) UnreadMessages(dto domain.Messages) {
	client := p.pusherConfig
	for _, unreadMessage := range dto.Items {
		data := map[string]int64{"NewMessageFrom": unreadMessage.SenderId}
		err := client.Trigger(strconv.FormatInt(unreadMessage.RecipientId, 10), "NewMessage", data)
		if err != nil {
			fmt.Println("error in trigger", strconv.FormatInt(unreadMessage.RecipientId, 10))
			fmt.Println(err.Error())
		}
	}
}
