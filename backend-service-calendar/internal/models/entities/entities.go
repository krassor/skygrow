package entities

import (
	"time"
)

type BaseModel struct {
	Id         string `gorm:"type:uuid;primary_key"`
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}

type Calendar struct {
	BaseModel
	OwnerId     string `gorm:"column:ownerId"`
	Description string `gorm:"column:description"`
	Etag        string `gorm:"column:etag"`
	Summary     string `gorm:"column:summary"`
	TimeZone    string `gorm:"column:timeZone"`
}

type CalendarEvent struct {
	BaseModel
	ConferenceId string    `gorm:"column:conferenceId"`
	RequestId    string    `gorm:"column:requestId"`
	Description  string    `gorm:"column:description"`
	Start        time.Time `gorm:"column:start"`
	End          time.Time `gorm:"column:end"`
	Status       string    `gorm:"column:status"`
	Summary      string    `gorm:"column:summary"`
	Transparency string    `gorm:"column:transparency"`
}

type Attendee struct {
	BaseModel
	UserId      string `gorm:"column:userId"`
	DisplayName string `gorm:"column:displayName"`
	Email       string `gorm:"column:email"`
}
