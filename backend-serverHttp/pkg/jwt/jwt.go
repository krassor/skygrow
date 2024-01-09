package jwt

import (
	"fmt"
	"time"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"

	"github.com/golang-jwt/jwt/v5"
)

// NewToken creates new JWT token for given user and app.
func NewToken(user entities.User, duration time.Duration, role string, secret string) (string, error) {
	op := "pkg.jwt.NewToken()"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["firstName"] = user.FirstName
	claims["secondName"] = user.SecondName
	claims["role"] = role
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return tokenString, nil
}
