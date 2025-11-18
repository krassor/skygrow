package handlers

import (
	"app/main.go/internal/config"
	"app/main.go/internal/models/domain"
	"app/main.go/internal/models/dto"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"

	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	//"github.com/yuin/goldmark"
)

var (
	ErrInternalServer = errors.New("internal server error")
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

type QuestionnaireHandler struct {
	LLMService LLMService
	cfg        *config.Config
	log        *slog.Logger
}

func NewQuestionnaireHandler(
	log *slog.Logger,
	cfg *config.Config,
	LLMService LLMService,
) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		LLMService: LLMService,
		log:        log,
		cfg:        cfg,
	}
}

func (h *QuestionnaireHandler) AdultCreate(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.AdultCreate()"
	requestID := uuid.New()
	log := h.log
	log.With(
		slog.String("op", op),
		slog.String("requestID", requestID.String()),
	)

	questionnaireDto := dto.AdultQuestionnaireDto{}

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

	log.Info(
		"start http handler Create",
	)

	user := domain.User{
		Email: questionnaireDto.User.Email,
		Name:  questionnaireDto.User.Name,
	}

	_, err = h.LLMService.AddJob(requestID, h.splitAdultQuestionnaire(&questionnaireDto), user, "ADULT")

	if err != nil {
		h.err(log, err, fmt.Errorf("Internal server error"), w, http.StatusBadRequest)
		return
	}

	err = utils.Json(w, http.StatusOK, dto.AdultResponseQuestionnaireDto{RequestID: requestID.String()})
	if err != nil {
		log.Error("error encode response to json", sl.Err(err))
	}

}

func (h *QuestionnaireHandler) SchoolchildCreate(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.QuestionnaireHandler.SchoolchildCreate()"
	requestID := uuid.New()
	log := h.log
	log.With(
		slog.String("op", op),
		slog.String("requestID", requestID.String()),
	)

	questionnaireDto := dto.SchoolchildQuestionnaireDto{}

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

	log.Info(
		"start http handler Create",
	)

	user := domain.User{
		Email: questionnaireDto.User.Email,
		Name:  questionnaireDto.User.Name,
	}

	_, err = h.LLMService.AddJob(requestID, h.splitSchoolchildQuestionnaire(&questionnaireDto), user, "SCHOOLCHILD")

	if err != nil {
		h.err(log, err, fmt.Errorf("Internal server error"), w, http.StatusBadRequest)
		return
	}

	err = utils.Json(w, http.StatusOK, dto.SchoolchildResponseQuestionnaireDto{RequestID: requestID.String()})
	if err != nil {
		log.Error("error encode response to json", sl.Err(err))
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

func (h *QuestionnaireHandler) splitAdultQuestionnaire(questionnaire *dto.AdultQuestionnaireDto) string {
	header := "Результат опроса пользователя:\n"
	result := fmt.Sprintf("%s\n", header) +
		splitAdultTest("Тест: Интересы (RIASEC)", questionnaire.RIASEC) +
		splitAdultTest("Тест: Объекты деятельности (Климов)", questionnaire.ObjectsOfActivityKlimov) +
		splitAdultTest("Тест: Личностные качества", questionnaire.PersonalQualities) +
		splitAdultTest("Тест: Ценности", questionnaire.Values)

	return result
}

func (h *QuestionnaireHandler) splitSchoolchildQuestionnaire(questionnaire *dto.SchoolchildQuestionnaireDto) string {
	header := "Результат опроса пользователя:\n"
	result := fmt.Sprintf("%s\n", header) +
		splitSchoolchildTest("Тест: Интересы (RIASEC)", questionnaire.RIASEC) +
		splitSchoolchildTest("Тест: Объекты деятельности (Климов)", questionnaire.ObjectsOfActivityKlimov) +
		splitSchoolchildTest("Тест: Личностные качества", questionnaire.PersonalQualities) +
		splitSchoolchildTest("Тест: Ценности", questionnaire.Values)

	return result
}

func splitAdultTest(testHeader string, answers []dto.AdultQuestionAnswer) string {
	result := fmt.Sprintf("%s\n", testHeader)
	for _, answer := range answers {
		result += fmt.Sprintf("Вопрос %d: %s\n", answer.Number, answer.Question) +
			fmt.Sprintf("Ответ: %s\n\n", answer.Answer)
	}
	return result
}

func splitSchoolchildTest(testHeader string, answers []dto.SchoolchildQuestionAnswer) string {
	result := fmt.Sprintf("%s\n", testHeader)
	for _, answer := range answers {
		result += fmt.Sprintf("Вопрос %d: %s\n", answer.Number, answer.Question) +
			fmt.Sprintf("Ответ: %s\n\n", answer.Answer)
	}
	return result
}

// func mdToHTML(md string) (string, error) {
// 	var buf bytes.Buffer
// 	if err := goldmark.Convert([]byte(md), &buf); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }
