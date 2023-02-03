package domain

import "time"

type Image struct {
	Id          int64
	SenderId    int64
	ObjId       int64
	ChatId      string
	ObjType     string
	Name        string
	RecipientId int64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (i Image) GetrecipientId() int64 {
	return i.RecipientId
}

func (i Image) GetUserId() int64 {
	return i.SenderId
}
