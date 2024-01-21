package dto

type RequestUserSignUpDto struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type ResponseUserSignUpDto struct {
	AccessToken string `json:"accessToken"`
}

type RequestUserSignInDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseUserSignInDto struct {
	AccessToken string `json:"accessToken"`
}
