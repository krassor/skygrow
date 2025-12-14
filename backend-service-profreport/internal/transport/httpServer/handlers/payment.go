package handlers

import (
	"app/main.go/internal/models/domain"
	"app/main.go/internal/models/dto"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// Payment обрабатывает уведомления от CloudPayments о успешной оплате
func (h *QuestionnaireHandler) Payment(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.Payment()"
	log := h.log.With(
		slog.String("op", op),
	)

	rBuf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(rBuf)
	if err != nil {
		h.err(log, err, fmt.Errorf("cannot read request body"), w, http.StatusBadRequest)
		return
	}

	log.Info(
		"payment notification received",
		slog.String("request body", string(rBuf)),
	)

	// Декодируем запрос от CloudPayments
	payRequestDto := dto.PayRequestDto{}
	err = json.NewDecoder(r.Body).Decode(&payRequestDto)
	if err != nil {
		h.err(log, err, fmt.Errorf("cannot decode json"), w, http.StatusBadRequest)
		return
	}

	log.Info(
		"received payment notification",
		slog.Int64("transaction_id", payRequestDto.TransactionId),
		slog.String("invoice_id", payRequestDto.InvoiceId),
		slog.String("status", payRequestDto.Status),
	)

	// Проверяем статус платежа
	if payRequestDto.Status != "Completed" && payRequestDto.Status != "Authorized" {
		h.err(log, fmt.Errorf("invalid payment status: %s", payRequestDto.Status), fmt.Errorf("invalid payment status"), w, http.StatusBadRequest)
		return
	}

	// InvoiceId должен содержать UUID опросника
	if payRequestDto.InvoiceId == "" {
		h.err(log, fmt.Errorf("InvoiceId is empty"), fmt.Errorf("InvoiceId is empty"), w, http.StatusBadRequest)
		return
	}

	requestID, err := uuid.Parse(payRequestDto.InvoiceId)
	if err != nil {
		h.err(log, err, fmt.Errorf("invalid InvoiceId format"), w, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Обновить запись в БД (установить PaymentSuccess=true, PaymentID=TransactionId)
	err = h.Repository.UpdatePaymentStatus(ctx, requestID, payRequestDto.TransactionId, true)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to update payment status"), w, http.StatusInternalServerError)
		return
	}

	// 2. Получить данные опросника из БД
	questionnaire, err := h.Repository.GetQuestionnaire(ctx, requestID)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to get questionnaire"), w, http.StatusInternalServerError)
		return
	}

	// 3. Получить данные пользователя из БД
	user, err := h.Repository.GetUser(ctx, questionnaire.UserID)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to get user"), w, http.StatusInternalServerError)
		return
	}

	// 4. Подготовить текст опросника на основе типа и запустить обработку через LLMService
	var questionnaireText string
	switch questionnaire.QuestionnaireType {
	case "ADULT":
		questionnaireText = buildAdultQuestionnaireText(&questionnaire.Answers)
	case "SCHOOLCHILD":
		questionnaireText = buildSchoolchildQuestionnaireText(&questionnaire.Answers)
	default:
		h.err(log, fmt.Errorf("unknown questionnaire type: %s", questionnaire.QuestionnaireType), fmt.Errorf("invalid questionnaire type"), w, http.StatusInternalServerError)
		return
	}

	// Создаем domain.User для LLM сервиса
	domainUser := domain.User{
		Email: user.Email,
		Name:  user.Name,
	}

	// Запускаем обработку через LLM
	_, err = h.LLMService.AddJob(requestID, questionnaireText, domainUser, questionnaire.QuestionnaireType)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to add LLM job"), w, http.StatusInternalServerError)
		return
	}

	log.Info("payment processed successfully",
		slog.String("request_id", requestID.String()),
		slog.String("questionnaire_type", questionnaire.QuestionnaireType),
		slog.Int64("transaction_id", payRequestDto.TransactionId),
	)

	// Возвращаем успешный ответ
	err = utils.Json(w, http.StatusOK, map[string]string{"code": "0"})
	if err != nil {
		log.Error("error encode response to json", sl.Err(err))
	}
}
