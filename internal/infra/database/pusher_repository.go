package database

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/resources"
	"fmt"
	"github.com/pusher/pusher-http-go/v5"
	"strconv"
)

//const pusherTableName = "pushers"

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

	data := map[string]string{"message from pusherRepository": fmt.Sprintf("add to contacts new user %s ", dto.Name)}
	err := client.Trigger("my-channel", "my-event", data)
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
	//fmt.Println(dto)
	client := p.pusherConfig

	//var unreadChats domain.Messages
	for _, unreadMessage := range dto.Items {
		data := map[string]int64{"NewMessageFrom": unreadMessage.SenderId}
		err := client.Trigger(strconv.FormatInt(unreadMessage.RecipientId, 10), "NewMessage", data)
		if err != nil {
			fmt.Println("error in trigger", strconv.FormatInt(unreadMessage.RecipientId, 10))
			fmt.Println(err.Error())
		}

		//
		//	unreadChat, err := s.messagesRepo.FindAllChatsWithUnreadMsg(chat.ChatId, id)
		//	if err != nil {
		//		log.Printf("messagesService:.messagesRepo.FindAllChatsWithUnreadMsg %s", err)
		//		return domain.Messages{}, err
		//	}
		//	unreadChats.Items = append(unreadChats.Items, unreadChat)
		//}
		//client := p.pusherConfig
		//var message []MessagesDto
		//
		//for _, c := range dto. {
		//	var dto MessageDto
		//	messageDto := dto.DomainToDto(c)
		//	result = append(result, messageDto)
		//}
		//
		//data := map[string]int64{"NewMessageFrom": dto.SenderId}
		//err := client.Trigger(strconv.FormatInt(dto.RecipientId, 10), "NewMessage", data)
		//if err != nil {
		//	fmt.Println("error in trigger", strconv.FormatInt(dto.RecipientId, 10))
		//	fmt.Println(err.Error())
		//}

	}
}
