package dto

type AdultQuestionAnswer struct {
	Number   int    `json:"number"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type AdultQuestionnaireDto struct {
	User                    UserDto               `json:"user"`
	Values                  []AdultQuestionAnswer `json:"values"`
	PersonalQualities       []AdultQuestionAnswer `json:"personalQualities"`
	ObjectsOfActivityKlimov []AdultQuestionAnswer `json:"objectsOfActivityKlimov"`
	RIASEC                  []AdultQuestionAnswer `json:"RIASEC"`
}

type AdultResponseQuestionnaireDto struct {
	RequestID string `json:"requestID"`
}

type SchoolchildQuestionAnswer struct {
	Number   int    `json:"number"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type SchoolchildQuestionnaireDto struct {
	User                    UserDto                     `json:"user"`
	Values                  []SchoolchildQuestionAnswer `json:"values"`
	PersonalQualities       []SchoolchildQuestionAnswer `json:"personalQualities"`
	ObjectsOfActivityKlimov []SchoolchildQuestionAnswer `json:"objectsOfActivityKlimov"`
	RIASEC                  []SchoolchildQuestionAnswer `json:"RIASEC"`
}

type SchoolchildResponseQuestionnaireDto struct {
	RequestID string `json:"requestID"`
}
