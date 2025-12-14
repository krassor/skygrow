package handlers

import (
	"app/main.go/internal/config"
	"app/main.go/internal/models/dto"
	"app/main.go/internal/models/repositories"
	"app/main.go/internal/utils"
	"app/main.go/internal/utils/logger/sl"

	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type QuestionnaireHandler struct {
	Repository Repository
	LLMService LLMService
	cfg        *config.Config
	log        *slog.Logger
}

func NewQuestionnaireHandler(
	log *slog.Logger,
	cfg *config.Config,
	LLMService LLMService,
	repo Repository,
) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		Repository: repo,
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

	err = checkAdultQuestionnaire(&questionnaireDto)
	if err != nil {
		h.err(log, err, err, w, http.StatusBadRequest)
		return
	}

	log.Info(
		"start http handler Create",
	)

	ctx := r.Context()

	// создаем или находим пользователя в БД
	dbUser := repositories.User{
		Name:  questionnaireDto.User.Name,
		Email: questionnaireDto.User.Email,
	}

	dbUser, err = h.Repository.FindOrCreateUser(ctx, dbUser)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to save user"), w, http.StatusInternalServerError)
		return
	}

	// создаем запись опросника в БД
	dbQuestionnaire := MapAdultQuestionnaireToRepository(&questionnaireDto, requestID, dbUser.ID)

	_, err = h.Repository.CreateQuestionnaire(ctx, dbQuestionnaire)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to save questionnaire"), w, http.StatusInternalServerError)
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

	err = checkSchoolchildQuestionnaire(&questionnaireDto)
	if err != nil {
		h.err(log, err, err, w, http.StatusBadRequest)
		return
	}

	log.Info(
		"start http handler Create",
	)

	ctx := r.Context()

	// создаем или находим пользователя в БД
	dbUser := repositories.User{
		Name:  questionnaireDto.User.Name,
		Email: questionnaireDto.User.Email,
	}

	dbUser, err = h.Repository.FindOrCreateUser(ctx, dbUser)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to save user"), w, http.StatusInternalServerError)
		return
	}

	// создаем запись опросника в БД
	dbQuestionnaire := MapSchoolchildQuestionnaireToRepository(&questionnaireDto, requestID, dbUser.ID)

	_, err = h.Repository.CreateQuestionnaire(ctx, dbQuestionnaire)
	if err != nil {
		h.err(log, err, fmt.Errorf("failed to save questionnaire"), w, http.StatusInternalServerError)
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

// func (h *QuestionnaireHandler) splitAdultQuestionnaire(questionnaire *dto.AdultQuestionnaireDto) string {
// 	header := "Результат опроса пользователя:\n"
// 	result := fmt.Sprintf("%s\n", header) +
// 		splitAdultTest("Тест: Интересы (RIASEC)", questionnaire.RIASEC) +
// 		splitAdultTest("Тест: Объекты деятельности (Климов)", questionnaire.ObjectsOfActivityKlimov) +
// 		splitAdultTest("Тест: Личностные качества", questionnaire.PersonalQualities) +
// 		splitAdultTest("Тест: Ценности", questionnaire.Values)

// 	return result
// }

// func (h *QuestionnaireHandler) splitSchoolchildQuestionnaire(questionnaire *dto.SchoolchildQuestionnaireDto) string {
// 	header := "Результат опроса пользователя:\n"
// 	result := fmt.Sprintf("%s\n", header) +
// 		splitSchoolchildTest("Тест: Интересы (RIASEC)", questionnaire.RIASEC) +
// 		splitSchoolchildTest("Тест: Объекты деятельности (Климов)", questionnaire.ObjectsOfActivityKlimov) +
// 		splitSchoolchildTest("Тест: Личностные качества", questionnaire.PersonalQualities) +
// 		splitSchoolchildTest("Тест: Ценности", questionnaire.Values)

// 	return result
// }

// func splitAdultTest(testHeader string, answers []dto.AdultQuestionAnswer) string {
// 	result := fmt.Sprintf("%s\n", testHeader)
// 	for _, answer := range answers {
// 		result += fmt.Sprintf("Вопрос %d: %s\n", answer.Number, answer.Question) +
// 			fmt.Sprintf("Ответ: %s\n\n", answer.Answer)
// 	}
// 	return result
// }

// func splitSchoolchildTest(testHeader string, answers []dto.SchoolchildQuestionAnswer) string {
// 	result := fmt.Sprintf("%s\n", testHeader)
// 	for _, answer := range answers {
// 		result += fmt.Sprintf("Вопрос %d: %s\n", answer.Number, answer.Question) +
// 			fmt.Sprintf("Ответ: %s\n\n", answer.Answer)
// 	}
// 	return result
// }

func checkAdultQuestionnaire(questionnaireDto *dto.AdultQuestionnaireDto) error {
	if questionnaireDto.User.Email == "" {
		return fmt.Errorf("email is empty")
	}
	if !utils.IsEmailValid(questionnaireDto.User.Email) {
		return fmt.Errorf("email is wrong: %s", questionnaireDto.User.Email)
	}
	if len(questionnaireDto.RIASEC) == 0 {
		return fmt.Errorf("riasec test is empty")
	}
	if len(questionnaireDto.ObjectsOfActivityKlimov) == 0 {
		return fmt.Errorf("objectsOfActivityKlimov test is empty")
	}
	if len(questionnaireDto.PersonalQualities) == 0 {
		return fmt.Errorf("personalQualities test is empty")
	}
	if len(questionnaireDto.Values) == 0 {
		return fmt.Errorf("values test is empty")
	}
	return nil
}

func checkSchoolchildQuestionnaire(questionnaireDto *dto.SchoolchildQuestionnaireDto) error {
	if questionnaireDto.User.Email == "" {
		return fmt.Errorf("email is empty")
	}
	if !utils.IsEmailValid(questionnaireDto.User.Email) {
		return fmt.Errorf("email is wrong: %s", questionnaireDto.User.Email)
	}
	if len(questionnaireDto.RIASEC) == 0 {
		return fmt.Errorf("riasec test is empty")
	}
	if len(questionnaireDto.ObjectsOfActivityKlimov) == 0 {
		return fmt.Errorf("objectsOfActivityKlimov test is empty")
	}
	if len(questionnaireDto.PersonalQualities) == 0 {
		return fmt.Errorf("personalQualities test is empty")
	}
	if len(questionnaireDto.Values) == 0 {
		return fmt.Errorf("values test is empty")
	}
	return nil
}
