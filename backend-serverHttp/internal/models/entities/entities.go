package entities

import (
	"gorm.io/gorm"
)

type Mentor struct {
	gorm.Model
	FirstName    string   `gorm:"column:firstName"`
	SecondName   string   `gorm:"column:secondName"`
	Phone        string   `gorm:"column:phone"`
	Email        string   `gorm:"column:email"`
	Competencies []string `gorm:"type:string[]"`
}

type User struct {
	gorm.Model
	FirstName    string   `gorm:"column:firstName"`
	SecondName   string   `gorm:"column:secondName"`
	Phone        string   `gorm:"column:phone"`
	Email        string   `gorm:"column:email"`
}

type BookOrder struct {
	gorm.Model
	FirstName  string `gorm:"column:firstName"`
	SecondName string `gorm:"column:secondName"`
	Phone      string `gorm:"column:phone"`
	Email      string `gorm:"column:email"`
	MentorID   uint   `gorm:"column:mentorID"`
}

type Subscriber struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	ChatID   int64  `gorm:"column:chatid"`
	IsActive bool   `gorm:"column:isActive"`
}
