package jwt

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-auth/internal/models/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserData struct {
	Uid        uuid.UUID `json:"uid"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Role       string    `json:"role"`
}
type MyClaims struct {
	jwt.RegisteredClaims
	Data UserData `json:"data"`
}

// NewToken creates new JWT token for given user and app.
func NewToken(user domain.User, duration time.Duration, role string, secret string) (string, error) {
	op := "pkg.jwt.NewToken()"

	token := jwt.New(jwt.SigningMethodHS256)

	userData := UserData{
		Uid:        user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Role:       role,
	}
	claims := token.Claims.(jwt.MapClaims)

	claims["data"] = userData
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return tokenString, nil
}
