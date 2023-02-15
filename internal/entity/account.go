package entity

import (
	"github.com/google/uuid"
)

type Account struct{
	ID uuid.UUID
	Login string
	Password string
}