package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"app/main.go/internal/models/repositories"

	"github.com/google/uuid"
)

func (r *Repository) FindQuestionnaireByID(ctx context.Context, id uuid.UUID) (repositories.Questionnaire, error) {
	var questionnaire repositories.Questionnaire
	query := `SELECT id, user_id, payment_id, payment_success, created_at, updated_at 
	          FROM questionnaires WHERE id = $1 LIMIT 1`

	err := r.DB.GetContext(ctx, &questionnaire, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repositories.Questionnaire{}, fmt.Errorf("questionnaire not found with id: %s", id)
		}
		return repositories.Questionnaire{}, fmt.Errorf("error in FindQuestionnaireByID(): %w", err)
	}

	return questionnaire, nil
}

func (r *Repository) FindQuestionnairesByUserID(ctx context.Context, userID uuid.UUID) ([]repositories.Questionnaire, error) {
	var questionnaires []repositories.Questionnaire
	query := `SELECT id, user_id, payment_id, payment_success, created_at, updated_at 
	          FROM questionnaires WHERE user_id = $1 ORDER BY created_at DESC`

	err := r.DB.SelectContext(ctx, &questionnaires, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error in FindQuestionnairesByUserID(): %w", err)
	}

	return questionnaires, nil
}

func (r *Repository) CreateQuestionnaire(ctx context.Context, questionnaire repositories.Questionnaire) (repositories.Questionnaire, error) {
	op := "repository.CreateQuestionnaire()"

	// генерация ID только если он не задан
	if questionnaire.ID == uuid.Nil {
		questionnaire.ID = uuid.New()
	}

	insertQuery := `INSERT INTO questionnaires (id, user_id, payment_id, payment_success, questionnaire_type) 
	                VALUES ($1, $2, $3, $4, $5)`

	_, err := r.DB.ExecContext(ctx, insertQuery,
		questionnaire.ID,
		questionnaire.UserID,
		questionnaire.PaymentID,
		questionnaire.PaymentSuccess,
		questionnaire.QuestionnaireType,
	)
	if err != nil {
		return repositories.Questionnaire{}, fmt.Errorf("%s: %w", op, err)
	}

	return questionnaire, nil
}

func (r *Repository) UpdateQuestionnaire(ctx context.Context, questionnaire repositories.Questionnaire) (repositories.Questionnaire, error) {
	updateQuery := `UPDATE questionnaires 
	                SET user_id = $1, payment_id = $2, payment_success = $3 
	                WHERE id = $4`

	result, err := r.DB.ExecContext(ctx, updateQuery,
		questionnaire.UserID,
		questionnaire.PaymentID,
		questionnaire.PaymentSuccess,
		questionnaire.ID,
	)
	if err != nil {
		return repositories.Questionnaire{}, fmt.Errorf("error in UpdateQuestionnaire(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return repositories.Questionnaire{}, fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repositories.Questionnaire{}, fmt.Errorf("questionnaire with id %s not found", questionnaire.ID)
	}

	return questionnaire, nil
}

func (r *Repository) DeleteQuestionnaire(ctx context.Context, id uuid.UUID) error {
	deleteQuery := `DELETE FROM questionnaires WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("error in DeleteQuestionnaire(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("questionnaire with id %s not found", id)
	}

	return nil
}

func (r *Repository) UpdatePaymentStatus(ctx context.Context, questionnaireID uuid.UUID, paymentSuccess bool) error {
	updateQuery := `UPDATE questionnaires 
	                SET payment_success = $1 
	                WHERE id = $2`

	result, err := r.DB.ExecContext(ctx, updateQuery, paymentSuccess, questionnaireID)
	if err != nil {
		return fmt.Errorf("error in UpdatePaymentStatus(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("questionnaire with id %s not found", questionnaireID)
	}

	return nil
}
