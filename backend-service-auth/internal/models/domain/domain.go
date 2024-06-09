package entities

import (
	"github.com/google/uuid"
	"time"
)

type BaseModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}

type User struct {
	BaseModel
	FirstName       string `gorm:"column:firstName"`
	SecondName      string `gorm:"column:secondName"`
	Phone           string `gorm:"column:phone"`
	Email           string `gorm:"column:email"`
	Hashed_password string `gorm:"column:hashedPassword"`
	AccessToken     string `gorm:"column:accessToken"`
}
