package entities

import (
	"time"
)

type BaseModel struct {
	ID         string `gorm:"type:uuid;primary_key"`
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}

type Mentor struct {
	BaseModel
	FirstName    string   `gorm:"column:firstName"`
	SecondName   string   `gorm:"column:secondName"`
	Phone        string   `gorm:"column:phone"`
	Email        string   `gorm:"column:email"`
	Competencies []string `gorm:"type:text[];column:competencies"`
}

type User struct {
	BaseModel
	FirstName  string `gorm:"column:firstName"`
	SecondName string `gorm:"column:secondName"`
	Phone      string `gorm:"column:phone"`
	Email      string `gorm:"column:email"`
}

// TODO: will be deprecated. BookOrder must be aggregate
type BookOrder struct {
	BaseModel
	FirstName          string `gorm:"column:firstName"`
	SecondName         string `gorm:"column:secondName"`
	Phone              string `gorm:"column:phone"`
	Email              string `gorm:"column:email"`
	MentorID           string `gorm:"column:mentorID"`
	ProblemDescription string `gorm:"type:string;column:problemDescription"`
}

type Subscriber struct {
	Name     string `gorm:"column:name"`
	ChatID   int64  `gorm:"column:chatid;primary_key"`
	IsActive bool   `gorm:"column:isActive"`
}
