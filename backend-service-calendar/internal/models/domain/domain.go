package domain

import (
	"github.com/google/uuid"
	"time"
)

type CalendarStatusType string

const (
	InProcess CalendarStatusType = "In process"
	Created   CalendarStatusType = "Created"
	Done      CalendarStatusType = "Done"
)

func (st CalendarStatusType) String() string {
	//switch st {
	//case InProcess:
	//
	//	return "In process"
	//case Created:
	//	return "Created"
	//case Done:
	//	return "Done"
	//default:
	//	return "Unknown"
	//}
	return string(st)
}

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type CalendarUser struct {
	BaseModel
	FirstName  string
	SecondName string
	Email      string `gorm:"index;column:email"`
}

type Calendar struct {
	BaseModel
	CalendarOwnerId  uuid.UUID `gorm:"column:calendar_owner_id;index"`
	GoogleCalendarId string
	Description      string
	Etag             string
	Summary          string
	TimeZone         string
	Status           CalendarStatusType
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
