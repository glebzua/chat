package domain

import (
	"time"
)

type Contact struct {
	Id          int64
	UserId      int64
	ContactId   int64
	Activated   bool
	ChatId      string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
	Nickname    string
}

type Contacts struct {
	Items []Contact
	Total uint64
}
