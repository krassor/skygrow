package dto

import "time"

type GoogleAuthToken struct {
	AccessToken  string    `json:"access_token"`
	Expiry       time.Time `json:"expiry"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
}
