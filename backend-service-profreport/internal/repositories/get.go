package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"app/main.go/internal/models/repositories"

	"github.com/google/uuid"
)

// GetQuestionnaire получает опросник по ID вместе с ответами
func (r *Repository) GetQuestionnaire(ctx context.Context, id uuid.UUID) (repositories.Questionnaire, error) {
	var questionnaire repositories.Questionnaire
	var answersJSON []byte

	query := `SELECT id, user_id, payment_id, payment_success, amount, questionnaire_type, answers, created_at, updated_at 
	          FROM questionnaires WHERE id = $1 LIMIT 1`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&questionnaire.ID,
		&questionnaire.UserID,
		&questionnaire.PaymentID,
		&questionnaire.PaymentSuccess,
		&questionnaire.Amount,
		&questionnaire.QuestionnaireType,
		&answersJSON,
		&questionnaire.CreatedAt,
		&questionnaire.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return repositories.Questionnaire{}, fmt.Errorf("questionnaire not found with id: %s", id)
		}
		return repositories.Questionnaire{}, fmt.Errorf("error in GetQuestionnaire(): %w", err)
	}

	// Декодируем JSON с ответами
	if len(answersJSON) > 0 {
		err = json.Unmarshal(answersJSON, &questionnaire.Answers)
		if err != nil {
			return repositories.Questionnaire{}, fmt.Errorf("error unmarshaling answers: %w", err)
		}
	}

	return questionnaire, nil
}

// GetUser получает пользователя по ID
func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (repositories.User, error) {
	var user repositories.User
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`

	err := r.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repositories.User{}, fmt.Errorf("user not found with id: %s", id)
		}
		return repositories.User{}, fmt.Errorf("error in GetUser(): %w", err)
	}

	return user, nil
}
