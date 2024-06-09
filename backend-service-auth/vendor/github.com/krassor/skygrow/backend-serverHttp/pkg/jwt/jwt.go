package jwt

import (
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	Uid        uuid.UUID `json:"uid"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Role       string    `json:"role"`
}

// NewToken creates new JWT token for given user and app.
func NewToken(user entities.User, duration time.Duration, role string, secret string) (string, error) {
	op := "pkg.jwt.NewToken()"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(*MyClaims)
	claims.Uid = user.ID
	claims.Email = user.Email
	claims.FirstName = user.FirstName
	claims.SecondName = user.SecondName
	claims.Role = role
	claims.IssuedAt.Time = time.Now()
	claims.ExpiresAt.Time = time.Now().Add(duration)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return tokenString, nil
}
