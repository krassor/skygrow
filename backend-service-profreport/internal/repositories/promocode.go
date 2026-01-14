package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"app/main.go/internal/models/repositories"

	"github.com/google/uuid"
)

// GetPromoCodeByCode retrieves a promo code by its code string (case-sensitive)
func (r *Repository) GetPromoCodeByCode(ctx context.Context, code string) (repositories.PromoCode, error) {
	var promoCode repositories.PromoCode
	query := `SELECT id, code, questionnaire_type, final_price, currency, expires_at, is_active, created_at, updated_at 
	          FROM promo_codes WHERE code = $1 AND is_active = true LIMIT 1`

	err := r.DB.GetContext(ctx, &promoCode, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return repositories.PromoCode{}, fmt.Errorf("promo code not found: %s", code)
		}
		return repositories.PromoCode{}, fmt.Errorf("error in GetPromoCodeByCode(): %w", err)
	}

	return promoCode, nil
}

// GetTestPriceByType retrieves the test price by questionnaire type
func (r *Repository) GetTestPriceByType(ctx context.Context, questionnaireType string) (repositories.TestPrice, error) {
	var testPrice repositories.TestPrice
	query := `SELECT id, questionnaire_type, price, currency, created_at, updated_at 
	          FROM test_prices WHERE questionnaire_type = $1 LIMIT 1`

	err := r.DB.GetContext(ctx, &testPrice, query, questionnaireType)
	if err != nil {
		if err == sql.ErrNoRows {
			return repositories.TestPrice{}, fmt.Errorf("test price not found for type: %s", questionnaireType)
		}
		return repositories.TestPrice{}, fmt.Errorf("error in GetTestPriceByType(): %w", err)
	}

	return testPrice, nil
}

// UpdatePaymentStatusWithPromoCode updates payment status using promo code (payment_id = 0 for promo)
func (r *Repository) UpdatePaymentStatusWithPromoCode(ctx context.Context, questionnaireID uuid.UUID) error {
	updateQuery := `UPDATE questionnaires 
	                SET payment_success = true, payment_id = 0 
	                WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, updateQuery, questionnaireID)
	if err != nil {
		return fmt.Errorf("error in UpdatePaymentStatusWithPromoCode(): %w", err)
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
