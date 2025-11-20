package repositories

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	BaseModel
	Name  string `db:"name"`
	Email string `db:"email"`
}

type Questionnaire struct {
	BaseModel
	UserID            uuid.UUID `db:"user_id"`
	PaymentID         uuid.UUID `db:"payment_id"`
	PaymentSuccess    bool      `db:"payment_success"`
	QuestionnaireType string    `db:"questionnaire_type"`
}
