package handlers

import (
	"app/main.go/internal/models/domain"
	"app/main.go/internal/models/dto"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Payment обрабатывает уведомления от CloudPayments о успешной оплате
func (h *QuestionnaireHandler) Payment(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.Payment()"
	log := h.log.With(
		slog.String("op", op),
	)

	// Читаем тело запроса для логирования
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.err(log, err, fmt.Errorf("cannot read request body"), w, http.StatusBadRequest)
		return
	}

	log.Debug(
		"payment notification received",
		slog.String("content_type", r.Header.Get("Content-Type")),
		slog.String("body", string(bodyBytes)),
		slog.Int("content_length", len(bodyBytes)),
	)

	// Декодируем запрос от CloudPayments
	payRequestDto := dto.PayRequestDto{}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// Восстанавливаем тело запроса для ParseForm
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if err := r.ParseForm(); err != nil {
			h.err(log, err, fmt.Errorf("failed to parse form"), w, http.StatusBadRequest)
			return
		}

		// Маппим основные поля из формы
		if tId, err := strconv.ParseInt(r.FormValue("TransactionId"), 10, 64); err == nil {
			payRequestDto.TransactionId = tId
		}
		if amt, err := strconv.ParseFloat(r.FormValue("Amount"), 64); err == nil {
			payRequestDto.Amount = amt
		}
		payRequestDto.Currency = r.FormValue("Currency")
		payRequestDto.DateTime = r.FormValue("DateTime")
		payRequestDto.CardFirstSix = r.FormValue("CardFirstSix")
		payRequestDto.CardLastFour = r.FormValue("CardLastFour")
		payRequestDto.CardType = r.FormValue("CardType")
		payRequestDto.CardExpDate = r.FormValue("CardExpDate")
		if tm, err := strconv.Atoi(r.FormValue("TestMode")); err == nil {
			payRequestDto.TestMode = tm
		}
		payRequestDto.Status = r.FormValue("Status")
		payRequestDto.OperationType = r.FormValue("OperationType")
		payRequestDto.GatewayName = r.FormValue("GatewayName")
		payRequestDto.InvoiceId = r.FormValue("InvoiceId")
		payRequestDto.AccountId = r.FormValue("AccountId")
		payRequestDto.Name = r.FormValue("Name")
		payRequestDto.Email = r.FormValue("Email")
	} else {
		// Пытаемся распарсить как JSON
		err = json.Unmarshal(bodyBytes, &payRequestDto)
		if err != nil {
			log.Error(
				"failed to decode payment request",
				sl.Err(err),
				slog.String("body_preview", string(bodyBytes[:min(200, len(bodyBytes))])),
			)
			h.err(log, err, fmt.Errorf("cannot decode json"), w, http.StatusBadRequest)
			return
		}
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

	// 2. Запустить процесс генерации отчета
	err = h.startReportGeneration(ctx, log, requestID)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to start report generation"), w, http.StatusInternalServerError)
		return
	}

	log.Info("payment processed successfully",
		slog.String("request_id", requestID.String()),
		slog.Int64("transaction_id", payRequestDto.TransactionId),
	)

	// Возвращаем успешный ответ
	err = utils.Json(w, http.StatusOK, map[string]string{"code": "0"})
	if err != nil {
		log.Error("error encode response to json", sl.Err(err))
	}
}

// startReportGeneration запускает процесс генерации отчета
// Используется при успешной оплате и при применении бесплатного промокода
func (h *QuestionnaireHandler) startReportGeneration(ctx context.Context, log *slog.Logger, requestID uuid.UUID) error {
	// 1. Получить данные опросника из БД
	questionnaire, err := h.Repository.GetQuestionnaire(ctx, requestID)
	if err != nil {
		return fmt.Errorf("failed to get questionnaire: %w", err)
	}

	// 2. Получить данные пользователя из БД
	user, err := h.Repository.GetUser(ctx, questionnaire.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 3. Подготовить текст опросника на основе типа
	var questionnaireText string
	switch questionnaire.QuestionnaireType {
	case "ADULT":
		questionnaireText = buildAdultQuestionnaireText(&questionnaire.Answers)
	case "SCHOOLCHILD":
		questionnaireText = buildSchoolchildQuestionnaireText(&questionnaire.Answers)
	default:
		return fmt.Errorf("unknown questionnaire type: %s", questionnaire.QuestionnaireType)
	}

	// 4. Создаем domain.User для LLM сервиса
	domainUser := domain.User{
		Email: user.Email,
		Name:  user.Name,
	}

	// 5. Запускаем обработку через LLM
	_, err = h.LLMService.AddJob(requestID, questionnaireText, domainUser, questionnaire.QuestionnaireType)
	if err != nil {
		return fmt.Errorf("failed to add LLM job: %w", err)
	}

	log.Info("report generation started",
		slog.String("request_id", requestID.String()),
		slog.String("questionnaire_type", questionnaire.QuestionnaireType),
	)

	return nil
}
