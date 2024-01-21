package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"time"
)

func ParseToken(tokenString string, secret string) (domain.User, error) {
	op := "utils.jwt.ParseToken()"

	var validMethods = []string{jwt.SigningMethodHS256.Name}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s:%w", op, errors.New("unexpected signing method"))
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods(validMethods))

	if err != nil {
		return domain.User{}, fmt.Errorf("%s:%w", op, err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return domain.User{}, fmt.Errorf("%s:%w", op, errors.New("error claims type assertion"))
	}

	if claims.ExpiresAt.Time.After(time.Now()) {
		return domain.User{}, fmt.Errorf("%s:%w", op, errors.New("token expired"))
	}

	// validate the essential claims
	if !token.Valid {
		return domain.User{}, fmt.Errorf("%s:%w", op, errors.New("invalid token"))
	}

	user := domain.User{
		BaseModel:  domain.BaseModel{ID: claims.Uid},
		FirstName:  claims.FirstName,
		SecondName: claims.SecondName,
		Email:      claims.Email,
	}

	return user, nil
}
