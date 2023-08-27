package userServices

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/dto"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	//telegramBot "github.com/krassor/skygrow/backend-serverHttp/internal/telegram"
)

var (
	ErrEmailNotValid    = errors.New("email is not valid")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type UserRepository interface {
	// FindAllBookOrder(ctx context.Context) ([]entities.BookOrder, error)
	// CreateBookOrder(ctx context.Context, bookOrder entities.BookOrder) (entities.BookOrder, error)
	// UpdateBookOrder(ctx context.Context, bookOrder entities.BookOrder) (entities.BookOrder, error)
	// FindBookOrderById(ctx context.Context, id string) (entities.BookOrder, error)
	FindUserByEmail(ctx context.Context, email string) (entities.User, error)
	CreateNewUser(ctx context.Context, user entities.User) (entities.User, error)
}

type UserService struct {
	UserRepository UserRepository
}

func NewUser(r UserRepository) *UserService {
	return &UserService{
		UserRepository: r,
	}
}

func (s *UserService) SignUp(ctx context.Context, userDto dto.RequestUserSignUpDto) error {
	if !isEmailValid(userDto.Email) || userDto.Email == "" {
		return ErrEmailNotValid
	}

	if userDto.Password == "" {
		return ErrEmptyPassword
	}

	userEntity := entities.User{
		FirstName:       userDto.FirstName,
		SecondName:      userDto.SecondName,
		Phone:           userDto.Phone,
		Email:           userDto.Email,
		Hashed_password: userDto.Password,
	}

	timeNow := time.Now()
	userEntity.BaseModel = entities.BaseModel{
		ID:         uuid.NewString(),
		Created_at: timeNow,
		Updated_at: timeNow,
	}

	findUserEntity, err := s.UserRepository.FindUserByEmail(ctx, userEntity.Email)
	if err != nil {
		return fmt.Errorf("error find user service: %w", err)
	}

	if (findUserEntity != entities.User{}) {
		return ErrUserAlreadyExist
	}

	_, err = s.UserRepository.CreateNewUser(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("error sign up service: %w", err)
	}

	return nil

}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
