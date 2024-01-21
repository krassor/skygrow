package domain

import (
	"github.com/google/uuid"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type User struct {
	BaseModel
	FirstName  string
	SecondName string
	Email      string `gorm:"index"`
}

type Calendar struct {
	BaseModel
	CalendarOwnerId  uuid.UUID `gorm:"column:calendar_owner_id;index"`
	GoogleCalendarId string
	Description      string
	Etag             string
	Summary          string
	TimeZone         string
}

type CalendarEvent struct {
	BaseModel
	CalendarId   uuid.UUID
	ConferenceId uuid.UUID
	RequestId    uuid.UUID
	Description  string
	Start        time.Time
	End          time.Time
	Status       string
	Summary      string
	Transparency string
}

type GoogleAuthToken struct {
	BaseModel
	AccessToken  string
	Expiry       time.Time
	RefreshToken string
	TokenType    string
}

//
//type Attendee struct {
//	ID          uuid.UUID
//	UserId      uuid.UUID
//	DisplayName string
//	Email       string
//}
