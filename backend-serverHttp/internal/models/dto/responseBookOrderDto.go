package dto

type ResponseBookOrderDto struct {
	FirstName          string `json:"firstName"`
	SecondName         string `json:"secondName"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	MentorID           uint   `json:"mentorID"`
	BookOrderID        uint   `json:"bookOrderID"`
	ProblemDescription string `json:"problemDescription"`
}
