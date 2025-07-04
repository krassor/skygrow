package userServices

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-auth/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-auth/internal/models/dto"
	"github.com/krassor/skygrow/backend-service-auth/pkg/jwt"
	"regexp"
	"time"
	//telegramBot "github.com/krassor/skygrow/backend-serverHttp/internal/telegram"
)

var (
	ErrEmailNotValid    = errors.New("email is not valid")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
	ErrWrongPassword    = errors.New("wrong password")
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateNewUser(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
}

type UserService struct {
	UserRepository UserRepository
	SecretKey      string
}

func NewUser(r UserRepository, secretKey string) *UserService {
	return &UserService{
		UserRepository: r,
		SecretKey:      secretKey,
	}
}

func (s *UserService) SignUp(ctx context.Context, userDto dto.RequestUserSignUpDto) error {
	if !isEmailValid(userDto.Email) || userDto.Email == "" {
		return ErrEmailNotValid
	}

	if userDto.Password == "" {
		return ErrEmptyPassword
	}

	userEntity := domain.User{
		FirstName:       userDto.FirstName,
		SecondName:      userDto.SecondName,
		Phone:           userDto.Phone,
		Email:           userDto.Email,
		Hashed_password: userDto.Password,
	}

	timeNow := time.Now()
	userEntity.BaseModel = domain.BaseModel{
		ID:         uuid.New(),
		Created_at: timeNow,
		Updated_at: timeNow,
	}

	findUserEntity, err := s.UserRepository.FindUserByEmail(ctx, userEntity.Email)
	if err != nil {
		return fmt.Errorf("error find user service: %w", err)
	}

	if (findUserEntity != domain.User{}) {
		return ErrUserAlreadyExist
	}

	_, err = s.UserRepository.CreateNewUser(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("error sign up service: %w", err)
	}

	return nil

}

// SignIn return accessToken and error
func (s *UserService) SignIn(ctx context.Context, userDto dto.RequestUserSignInDto) (dto.ResponseUserSignInDto, error) {
	if !isEmailValid(userDto.Email) || userDto.Email == "" {
		return dto.ResponseUserSignInDto{}, ErrEmailNotValid
	}

	if userDto.Password == "" {
		return dto.ResponseUserSignInDto{}, ErrEmptyPassword
	}

	findUserEntity, err := s.UserRepository.FindUserByEmail(ctx, userDto.Email)
	if err != nil {
		return dto.ResponseUserSignInDto{}, fmt.Errorf("error find user. Service SignIn(): %w", err)
	}

	if (findUserEntity == domain.User{}) {
		return dto.ResponseUserSignInDto{}, ErrUserNotFound
	}

	if findUserEntity.Hashed_password != userDto.Password {
		return dto.ResponseUserSignInDto{}, ErrWrongPassword
	}
	accessToken, err := jwt.NewToken(findUserEntity, time.Hour*72, "user", s.SecretKey)

	if err != nil {
		return dto.ResponseUserSignInDto{}, fmt.Errorf("error generate access token. Service SignIn(): %w", err)
	}
	//end generate jwt token

	findUserEntity.AccessToken = accessToken
	findUserEntity.BaseModel.Updated_at = time.Now()

	_, err = s.UserRepository.UpdateUser(ctx, findUserEntity)
	if err != nil {
		return dto.ResponseUserSignInDto{}, fmt.Errorf("error update user. Service SignIn(): %w", err)
	}

	return dto.ResponseUserSignInDto{AccessToken: accessToken}, nil

}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
