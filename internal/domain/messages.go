package domain

import (
	"time"
)

type Message struct {
	Id          int64
	ChatId      string
	SenderId    int64
	RecipientId int64
	Message     string
	FileLoc     string
	Sended      bool
	Received    bool
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Messages struct {
	Items []Message
}
type Chat struct {
	Id          int64
	ChatId      string
	SenderId    int64
	RecipientId int64
	Message     string
	Sended      bool
	Received    bool
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Chats struct {
	Items []Chat
}
