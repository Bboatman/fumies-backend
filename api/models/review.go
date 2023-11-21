package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Title     string
	Body      string
	Status    string
	CreatedAt time.Time
}
