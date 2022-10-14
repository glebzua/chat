package domain

import (
	"time"
)

type Role string

const (
	ROLE_ADMIN Role = "admin"
	ROLE_USER  Role = "user"
)

type User struct {
	Id          int64
	Email       string
	Password    string
	Name        string
	PhoneNumber string
	Avatar      *string
	Activated   bool
	Role        Role
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Users struct {
	Items []User
	Total uint64
}

type ChangePassword struct {
	OldPassword string
	NewPassword string
}

func (u User) GetUserId() int64 {
	return u.Id
}

type ForcedDelete struct {
	Forced bool
}
