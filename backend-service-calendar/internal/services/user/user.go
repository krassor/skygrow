package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils"
	"regexp"

	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	//telegramBot "github.com/krassor/skygrow/backend-service-calendar/internal/telegram"
)

var (
	ErrEmailNotValid    = errors.New("email is not valid")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
	ErrWrongPassword    = errors.New("wrong password")
)

type userRepository interface {
	FindById(ctx context.Context, userId uuid.UUID) (domain.CalendarUser, error)
	FindByEmail(ctx context.Context, email string) (domain.CalendarUser, error)
	CreateUser(ctx context.Context, email string) (domain.CalendarUser, error)
}

type User struct {
	userRepository userRepository
}

func NewUser(r userRepository) *User {
	return &User{userRepository: r}
}

func (u *User) FindUserById(ctx context.Context, userId uuid.UUID) (domain.CalendarUser, error) {
	op := "CalendarService FindUserById()"
	user, err := u.userRepository.FindById(ctx, userId)
	if err != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, err)

	}
	return user, nil
}

func (u *User) FindUserByEmail(ctx context.Context, email string) (domain.CalendarUser, error) {
	op := "CalendarService FindUserByEmail()"
	if !utils.IsEmailValid(email) {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, ErrEmailNotValid)
	}
	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, err)

	}
	return user, nil
}

func (u *User) CreateUser(ctx context.Context, email string) (domain.CalendarUser, error) {
	op := "CalendarService CreateUser()"
	if !utils.IsEmailValid(email) {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, ErrEmailNotValid)
	}
	user, err := u.userRepository.CreateUser(ctx, email)
	if err != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s:%w", op, err)

	}
	return user, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
