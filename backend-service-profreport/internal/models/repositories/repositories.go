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
	PaymentID         int64     `db:"payment_id"`
	PaymentSuccess    bool      `db:"payment_success"`
	Amount            int       `db:"amount"`
	QuestionnaireType string    `db:"questionnaire_type"`
	Answers           Answers   `db:"answers"`
}

type QuestionAnswer struct {
	Number   int    `db:"number"`
	Question string `db:"question"`
	Answer   string `db:"answer"`
}

type Answers struct {
	Values                  []QuestionAnswer `db:"values"`
	PersonalQualities       []QuestionAnswer `db:"personal_qualities"`
	ObjectsOfActivityKlimov []QuestionAnswer `db:"objects_of_activity_klimov"`
	RIASEC                  []QuestionAnswer `db:"riasec"`
}
