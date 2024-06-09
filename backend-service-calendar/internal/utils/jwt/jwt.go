package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	ErrorClaimsTypeAssertion = errors.New("error claims type assertion")
	ErrorNoExpiresAtClaim    = errors.New("error exp claim not found in the jwt token")
)

func ParseToken(tokenString string, secret string) (domain.CalendarUser, error) {
	op := "utils.jwt.ParseToken()"

	var validMethods = []string{jwt.SigningMethodHS256.Name}

	//token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("%s:%w", op, errors.New("unexpected signing method"))
	//	}
	//	return []byte(secret), nil
	//}, jwt.WithValidMethods(validMethods))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s:%w", op, errors.New("unexpected signing method"))
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods(validMethods))

	if err != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, ErrorClaimsTypeAssertion)
	}

	// validate the essential claims
	if !token.Valid {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, errors.New("invalid token"))
	}

	userDataMap := claims["data"].(map[string]interface{})
	uid, err := uuid.Parse(userDataMap["uid"].(string))
	if err != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, err)
	}

	userData := UserData{
		Uid:        uid,
		Email:      userDataMap["email"].(string),
		FirstName:  userDataMap["first_name"].(string),
		SecondName: userDataMap["second_name"].(string),
		Role:       userDataMap["role"].(string),
	}

	et, err := claims.GetExpirationTime()
	if err != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, err)
	}
	if et == nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, ErrorNoExpiresAtClaim)
	}

	log.Debug().Msgf("%s : UserData: %v, et: %v, timeNow: %v", op, userData, et, time.Now())

	if et.Before(time.Now()) {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, errors.New("token expired"))
	}

	user := domain.CalendarUser{
		BaseModel:  domain.BaseModel{ID: userData.Uid},
		FirstName:  userData.FirstName,
		SecondName: userData.SecondName,
		Email:      userData.Email,
	}

	return user, nil
}
