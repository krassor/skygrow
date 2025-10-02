package handlers

import (
	"app/main.go/internal/mail"
	"app/main.go/internal/models/dto"
	"app/main.go/internal/openrouter"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// "app/main.go/internal/models/domain"
	// "app/main.go/internal/models/dto"
	// "app/main.go/internal/utils"
	// "app/main.go/internal/utils/logger/sl"
	"log/slog"
	"net/http"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type QuestionnaireHandler struct {
	openrouterService *openrouter.Openrouter
	mailerService     *mail.Mailer
	log               *slog.Logger
}

func NewQuestionnaireHandler(log *slog.Logger, openrouterService *openrouter.Openrouter, mailer *mail.Mailer) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		openrouterService: openrouterService,
		mailerService:     mailer,
		log:               log,
	}
}

func (h *QuestionnaireHandler) Create(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.Create()"
	log := h.log
	log.With(
		slog.String("op", op),
	)

	questionnaireDto := dto.QuestionnaireDto{}

	err := json.NewDecoder(r.Body).Decode(&questionnaireDto)

	if err != nil {
		h.err(log, err, fmt.Errorf("cannot decode json"), w, http.StatusBadRequest)
		return
	}

	if questionnaireDto.User.Email == "" {
		h.err(log, fmt.Errorf("email is empty"), fmt.Errorf("email is empty"), w, http.StatusBadRequest)
		return
	}
	if len(questionnaireDto.RIASEC) == 0 {
		h.err(log, fmt.Errorf("RIASEC test is empty"), fmt.Errorf("RIASEC test is empty"), w, http.StatusBadRequest)
		return
	}
	if len(questionnaireDto.ObjectsOfActivityKlimov) == 0 {
		h.err(log, fmt.Errorf("ObjectsOfActivityKlimov test is empty"), fmt.Errorf("ObjectsOfActivityKlimov test is empty"), w, http.StatusBadRequest)
	}
	if len(questionnaireDto.PersonalQualities) == 0 {
		h.err(log, fmt.Errorf("PersonalQualities test is empty"), fmt.Errorf("PersonalQualities test is empty"), w, http.StatusBadRequest)
		return
	}
	if len(questionnaireDto.Values) == 0 {
		h.err(log, fmt.Errorf("Values test is empty"), fmt.Errorf("Values test is empty"), w, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	response, err := h.openrouterService.CreateChatCompletion(ctx, h.splitQuestionnaire(&questionnaireDto))
	if err != nil {
		h.err(log, err, fmt.Errorf("internal server error"), w, http.StatusInternalServerError)
		return
	}

	mailBody := "Здравствуйте!\n По Вашему запросу был сгенерирован отчет\n" + response + "\nС уважением, команда profreport."

	err = h.mailerService.AddJob(questionnaireDto.User.Email, "Prof Report", mailBody)
	if err != nil {
		h.err(log, err, fmt.Errorf("internal server error"), w, http.StatusInternalServerError)
		return
	}

}

func (h *QuestionnaireHandler) err(log *slog.Logger, err error, httpErr error, w http.ResponseWriter, status int) {

	log.Error(
		"error",
		sl.Err(err),
	)
	httpError := utils.Err(w, status, httpErr)
	if httpError != nil {
		log.Error(
			"error sending http response",
			sl.Err(err),
		)
	}

}

func (h *QuestionnaireHandler) splitQuestionnaire(questionnaire *dto.QuestionnaireDto) string {
	header := "Результат опроса пользователя:\n"
	result := fmt.Sprintf("%s\n", header) +
		splitTest("Тест: Интересы (RIASEC)", questionnaire.RIASEC) +
		splitTest("Тест: Объекты деятельности (Климов)", questionnaire.ObjectsOfActivityKlimov) +
		splitTest("Тест: Личностные качества", questionnaire.PersonalQualities) +
		splitTest("Тест: Ценности", questionnaire.Values)

	return result
}

func splitTest(testHeader string, answers []dto.QuestionAnswer) string {
	result := fmt.Sprintf("%s\n", testHeader)
	for _, answer := range answers {
		result += fmt.Sprintf("Вопрос %d: %s\n", answer.Number, answer.Question) +
			fmt.Sprintf("Ответ: %s\n\n", answer.Answer)
	}
	return result
}
