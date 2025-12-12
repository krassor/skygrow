package handlers

import (
	"context"

	"app/main.go/internal/models/domain"
	"app/main.go/internal/models/repositories"

	"github.com/google/uuid"
)

// Параметры:
//   - requestID: уникальный идентификатор запроса (UUID).
//   - jobType: тип запроса: ADULT, SCHOOLCHILD
type LLMService interface {
	AddJob(
		reqyestID uuid.UUID,
		questionnaire string,
		user domain.User,
		jobType string,
	) (chan struct{}, error)
}

// Repository интерфейс для работы с БД
type Repository interface {
	FindOrCreateUser(ctx context.Context, user repositories.User) (repositories.User, error)
	CreateQuestionnaire(ctx context.Context, questionnaire repositories.Questionnaire) (repositories.Questionnaire, error)
	GetQuestionnaire(ctx context.Context, id uuid.UUID) (repositories.Questionnaire, error)
	GetUser(ctx context.Context, id uuid.UUID) (repositories.User, error)
	UpdatePaymentStatus(ctx context.Context, questionnaireID uuid.UUID, paymentID int64, paymentSuccess bool) error
}
