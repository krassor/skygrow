package handlers

import (
	"app/main.go/internal/models/dto"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ApplyPromoCode обрабатывает применение промокода
func (h *QuestionnaireHandler) ApplyPromoCode(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.ApplyPromoCode()"
	log := h.log.With(
		slog.String("op", op),
	)

	// Декодируем запрос
	var req dto.ApplyPromoCodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.err(log, err, fmt.Errorf("cannot decode json"), w, http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.Promocode == "" {
		h.errJson(log, fmt.Errorf("promocode is empty"), w, http.StatusBadRequest, "Промокод не указан")
		return
	}

	if req.RequestID == "" {
		h.errJson(log, fmt.Errorf("requestID is empty"), w, http.StatusBadRequest, "RequestID не указан")
		return
	}

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		h.errJson(log, err, w, http.StatusBadRequest, "Неверный формат RequestID")
		return
	}

	ctx := r.Context()

	// Получаем опросник для определения типа теста
	questionnaire, err := h.Repository.GetQuestionnaire(ctx, requestID)
	if err != nil {
		h.errJson(log, err, w, http.StatusNotFound, "Анкета не найдена")
		return
	}

	// Поиск промокода (регистрочувствительный)
	promoCode, err := h.Repository.GetPromoCodeByCode(ctx, req.Promocode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == fmt.Sprintf("promo code not found: %s", req.Promocode) {
			log.Debug("promo code not found", slog.String("code", req.Promocode))
			err = utils.Json(w, http.StatusBadRequest, dto.PromoCodeErrorResponse{Error: "Недействительный промокод"})
			if err != nil {
				log.Error("error encode response to json", sl.Err(err))
			}
			return
		}
		h.err(log, err, fmt.Errorf("failed to get promo code"), w, http.StatusInternalServerError)
		return
	}

	// Проверяем, что промокод соответствует типу теста
	if promoCode.QuestionnaireType != questionnaire.QuestionnaireType {
		log.Debug("promo code type mismatch",
			slog.String("promo_type", promoCode.QuestionnaireType),
			slog.String("questionnaire_type", questionnaire.QuestionnaireType),
		)
		err = utils.Json(w, http.StatusBadRequest, dto.PromoCodeErrorResponse{Error: "Недействительный промокод"})
		if err != nil {
			log.Error("error encode response to json", sl.Err(err))
		}
		return
	}

	// Проверяем срок действия промокода
	if time.Now().After(promoCode.ExpiresAt) {
		log.Debug("promo code expired",
			slog.String("code", req.Promocode),
			slog.Time("expires_at", promoCode.ExpiresAt),
		)
		err = utils.Json(w, http.StatusBadRequest, dto.PromoCodeErrorResponse{Error: "Время действия промокода истекло"})
		if err != nil {
			log.Error("error encode response to json", sl.Err(err))
		}
		return
	}

	// Получаем исходную стоимость теста
	testPrice, err := h.Repository.GetTestPriceByType(ctx, questionnaire.QuestionnaireType)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to get test price"), w, http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := dto.ApplyPromoCodeResponse{
		Promocode:     req.Promocode,
		OriginalPrice: testPrice.Price,
		FinalPrice:    promoCode.FinalPrice,
		Currency:      testPrice.Currency,
	}

	// Если конечная стоимость = 0, запускаем генерацию отчета
	if promoCode.FinalPrice == 0 {
		// Обновляем статус оплаты
		err = h.Repository.UpdatePaymentStatusWithPromoCode(ctx, requestID)
		if err != nil {
			h.err(log, err, fmt.Errorf("failed to update payment status"), w, http.StatusInternalServerError)
			return
		}

		// Запускаем генерацию отчета
		err = h.startReportGeneration(ctx, log, requestID)
		if err != nil {
			h.err(log, err, fmt.Errorf("failed to start report generation"), w, http.StatusInternalServerError)
			return
		}

		log.Info("free promo code applied, report generation started",
			slog.String("request_id", requestID.String()),
			slog.String("promo_code", req.Promocode),
		)
	} else {
		log.Info("promo code applied",
			slog.String("request_id", requestID.String()),
			slog.String("promo_code", req.Promocode),
			slog.Int("original_price", testPrice.Price),
			slog.Int("final_price", promoCode.FinalPrice),
		)
	}

	// Возвращаем успешный ответ
	err = utils.Json(w, http.StatusOK, response)
	if err != nil {
		log.Error("error encode response to json", sl.Err(err))
	}
}

// errJson отправляет JSON с ошибкой
func (h *QuestionnaireHandler) errJson(log *slog.Logger, err error, w http.ResponseWriter, status int, message string) {
	log.Error("error", sl.Err(err), slog.String("message", message))
	httpError := utils.Json(w, status, dto.PromoCodeErrorResponse{Error: message})
	if httpError != nil {
		log.Error("error sending http response", sl.Err(httpError))
	}
}
