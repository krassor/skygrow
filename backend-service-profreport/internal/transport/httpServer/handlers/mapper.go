package handlers

import (
	"app/main.go/internal/models/dto"
	"app/main.go/internal/models/repositories"

	"github.com/google/uuid"
)

// MapAdultQuestionnaireToRepository преобразует AdultQuestionnaireDto в repositories.Questionnaire
func MapAdultQuestionnaireToRepository(
	questionnaireDto *dto.AdultQuestionnaireDto,
	requestID uuid.UUID,
	userID uuid.UUID,
) repositories.Questionnaire {
	return repositories.Questionnaire{
		BaseModel: repositories.BaseModel{
			ID: requestID,
		},
		UserID:            userID,
		PaymentID:         0,
		PaymentSuccess:    false,
		Amount:            0,
		QuestionnaireType: "ADULT",
		Answers: repositories.Answers{
			Values:                  mapAdultQuestionAnswersToRepository(questionnaireDto.Values),
			PersonalQualities:       mapAdultQuestionAnswersToRepository(questionnaireDto.PersonalQualities),
			ObjectsOfActivityKlimov: mapAdultQuestionAnswersToRepository(questionnaireDto.ObjectsOfActivityKlimov),
			RIASEC:                  mapAdultQuestionAnswersToRepository(questionnaireDto.RIASEC),
		},
	}
}

// MapSchoolchildQuestionnaireToRepository преобразует SchoolchildQuestionnaireDto в repositories.Questionnaire
func MapSchoolchildQuestionnaireToRepository(
	questionnaireDto *dto.SchoolchildQuestionnaireDto,
	requestID uuid.UUID,
	userID uuid.UUID,
) repositories.Questionnaire {
	return repositories.Questionnaire{
		BaseModel: repositories.BaseModel{
			ID: requestID,
		},
		UserID:            userID,
		PaymentID:         0,
		PaymentSuccess:    false,
		Amount:            0,
		QuestionnaireType: "SCHOOLCHILD",
		Answers: repositories.Answers{
			Values:                  mapSchoolchildQuestionAnswersToRepository(questionnaireDto.Values),
			PersonalQualities:       mapSchoolchildQuestionAnswersToRepository(questionnaireDto.PersonalQualities),
			ObjectsOfActivityKlimov: mapSchoolchildQuestionAnswersToRepository(questionnaireDto.ObjectsOfActivityKlimov),
			RIASEC:                  mapSchoolchildQuestionAnswersToRepository(questionnaireDto.RIASEC),
		},
	}
}

// mapAdultQuestionAnswersToRepository преобразует массив AdultQuestionAnswer в массив repositories.QuestionAnswer
func mapAdultQuestionAnswersToRepository(answers []dto.AdultQuestionAnswer) []repositories.QuestionAnswer {
	result := make([]repositories.QuestionAnswer, len(answers))
	for i, answer := range answers {
		result[i] = repositories.QuestionAnswer{
			Number:   answer.Number,
			Question: answer.Question,
			Answer:   answer.Answer,
		}
	}
	return result
}

// mapSchoolchildQuestionAnswersToRepository преобразует массив SchoolchildQuestionAnswer в массив repositories.QuestionAnswer
func mapSchoolchildQuestionAnswersToRepository(answers []dto.SchoolchildQuestionAnswer) []repositories.QuestionAnswer {
	result := make([]repositories.QuestionAnswer, len(answers))
	for i, answer := range answers {
		result[i] = repositories.QuestionAnswer{
			Number:   answer.Number,
			Question: answer.Question,
			Answer:   answer.Answer,
		}
	}
	return result
}
