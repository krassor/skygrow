package handlers

import (
	"app/main.go/internal/config"
	"app/main.go/internal/models/dto"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type LLMService interface {
	CreateChatCompletion(
		ctx context.Context,
		logger *slog.Logger,
		requestId uuid.UUID,
		prompt string,
	) (string, error)
}

type MailService interface {
	AddJob(
		ID uuid.UUID,
		to string,
		subject string,
		body string,
	) error
}

type PdfService interface {
	AddJob(
		requestId uuid.UUID,
		inputMarkdown string,
	) (chan struct{}, error)
}

type QuestionnaireHandler struct {
	LLMService  LLMService
	MailService MailService
	PdfService  PdfService
	cfg         *config.Config
	log         *slog.Logger
}

func NewQuestionnaireHandler(
	log *slog.Logger,
	cfg *config.Config,
	LLMService LLMService,
	MailService MailService,
	PdfService PdfService,
) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		LLMService:  LLMService,
		MailService: MailService,
		PdfService:  PdfService,
		log:         log,
		cfg:         cfg,
	}
}

func (h *QuestionnaireHandler) Create(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.Create()"
	requestID := uuid.New()
	log := h.log
	log.With(
		slog.String("op", op),
		slog.String("requestID", requestID.String()),
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
	if !utils.IsEmailValid(questionnaireDto.User.Email) {
		h.err(log, fmt.Errorf("email is wrong"), fmt.Errorf("email is wrong: %s", questionnaireDto.User.Email), w, http.StatusBadRequest)
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

	responseSuccess := utils.Message(true, "")
	err = utils.Json(w, http.StatusOK, responseSuccess)
	if err != nil {
		log.Error("error encode response to json", sl.Err(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.BotConfig.AI.GetTimeout())
	defer cancel()

	response, err := h.LLMService.CreateChatCompletion(ctx, log, requestID, h.splitQuestionnaire(&questionnaireDto))
	if err != nil {
		sl.Err(err)
		return
	}

	chanDone, err := h.PdfService.AddJob(requestID, response)
	if err != nil {
		sl.Err(err)
		return
	}

	<-chanDone

	mailBody := "Здравствуйте, " + questionnaireDto.User.Name + "!\r\n" +
		"По Вашему запросу был сгенерирован отчет\r\n" +
		"Отчет прикреплен к письму во вложении.\r\n" +
		"\r\n\r\nС уважением, команда proffreport."

	mailBody, err = mdToHTML(mailBody)
	if err != nil {
		sl.Err(err)
		return
	}

	err = h.MailService.AddJob(requestID, questionnaireDto.User.Email, "Prof Report", mailBody)
	if err != nil {
		sl.Err(err)
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

func mdToHTML(md string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
