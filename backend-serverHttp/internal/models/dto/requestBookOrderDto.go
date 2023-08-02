package dto

type RequestBookOrderDto struct {
	FirstName          string `json:"firstName"`
	SecondName         string `json:"secondName"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	MentorID           string `json:"mentorID"`
	ProblemDescription string `json:"problemDescription"`
}
