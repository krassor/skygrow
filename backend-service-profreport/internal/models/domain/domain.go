package domain

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

type Payment struct {
	ID             uuid.UUID
	PaymentSuccess bool
}

type Questionnaire struct {
	ID      uuid.UUID
	User    User
	Payment Payment
	Type    string
}
