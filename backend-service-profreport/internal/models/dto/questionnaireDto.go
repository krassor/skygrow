package dto

type QuestionAnswer struct {
	Number   int    `json:"number"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type QuestionnaireDto struct {
	User                    UserDto          `json:"user"`
	Values                  []QuestionAnswer `json:"values"`
	PersonalQualities       []QuestionAnswer `json:"personalQualities"`
	ObjectsOfActivityKlimov []QuestionAnswer `json:"objectsOfActivityKlimov"`
	RIASEC                  []QuestionAnswer `json:"RIASEC"`
}
