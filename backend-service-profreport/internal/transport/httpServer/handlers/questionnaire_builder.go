package handlers

import (
	"app/main.go/internal/models/repositories"
	"fmt"
)

// buildAdultQuestionnaireText формирует текст опросника для взрослых из структуры Answers
func buildAdultQuestionnaireText(answers *repositories.Answers) string {
	header := "Результат опроса пользователя:\n"
	result := fmt.Sprintf("%s\n", header) +
		buildTestSection("Тест: Интересы (RIASEC)", answers.RIASEC) +
		buildTestSection("Тест: Объекты деятельности (Климов)", answers.ObjectsOfActivityKlimov) +
		buildTestSection("Тест: Личностные качества", answers.PersonalQualities) +
		buildTestSection("Тест: Ценности", answers.Values)

	return result
}

// buildSchoolchildQuestionnaireText формирует текст опросника для школьников из структуры Answers
func buildSchoolchildQuestionnaireText(answers *repositories.Answers) string {
	header := "Результат опроса пользователя:\n"
	result := fmt.Sprintf("%s\n", header) +
		buildTestSection("Тест: Интересы (RIASEC)", answers.RIASEC) +
		buildTestSection("Тест: Объекты деятельности (Климов)", answers.ObjectsOfActivityKlimov) +
		buildTestSection("Тест: Личностные качества", answers.PersonalQualities) +
		buildTestSection("Тест: Ценности", answers.Values)

	return result
}

// buildTestSection формирует секцию для одного теста
func buildTestSection(testHeader string, answers []repositories.QuestionAnswer) string {
	result := fmt.Sprintf("%s\n", testHeader)
	for _, answer := range answers {
		result += fmt.Sprintf("Вопрос %d: %s\n", answer.Number, answer.Question) +
			fmt.Sprintf("Ответ: %s\n\n", answer.Answer)
	}
	return result
}
