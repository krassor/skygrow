package repositories

import (
	"time"

	"github.com/google/uuid"
)

// TestPrice represents test type pricing
type TestPrice struct {
	ID                int       `db:"id"`
	QuestionnaireType string    `db:"questionnaire_type"`
	Price             int       `db:"price"`
	Currency          string    `db:"currency"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// PromoCode represents a promotional code
type PromoCode struct {
	ID                uuid.UUID `db:"id"`
	Code              string    `db:"code"`
	QuestionnaireType string    `db:"questionnaire_type"`
	FinalPrice        int       `db:"final_price"`
	Currency          string    `db:"currency"`
	ExpiresAt         time.Time `db:"expires_at"`
	IsActive          bool      `db:"is_active"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
