package domain

import (
	"github.com/google/uuid"
)

type Session struct {
	UserId int64
	UUID   uuid.UUID
}
