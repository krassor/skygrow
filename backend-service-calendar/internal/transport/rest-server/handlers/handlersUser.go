package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/krassor/skygrow/backend-service-calendar/internal/models/dto"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/userServices"
	"github.com/krassor/skygrow/backend-service-calendar/pkg/utils"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *userServices.UserService
}

func NewUserHandler(userService *userServices.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	requestUserSignUpDto := dto.RequestUserSignUpDto{}

	err := json.NewDecoder(r.Body).Decode(&requestUserSignUpDto)
	if err != nil {
		log.Error().Msgf("Error decode json in SignUp() handler: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Error().Msgf("Cannot sending error message to http client from SignUp(): %s", httpErr)
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = h.userService.SignUp(ctx, requestUserSignUpDto)

	if err != nil {
		if errors.Is(err, userServices.ErrEmailNotValid) {
			log.Error().Msgf("Error creating user: %v: %s", requestUserSignUpDto, err)
			httpErr := utils.Err(w, http.StatusBadRequest, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		if errors.Is(err, userServices.ErrEmptyPassword) {
			log.Error().Msgf("Error creating user: %v: %s", requestUserSignUpDto, err)
			httpErr := utils.Err(w, http.StatusBadRequest, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		if errors.Is(err, userServices.ErrUserAlreadyExist) {
			log.Error().Msgf("Error creating user: %v: %s", requestUserSignUpDto, err)
			httpErr := utils.Err(w, http.StatusBadRequest, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		log.Error().Msgf("Error creating user: %v: %s", requestUserSignUpDto, err)
		httpErr := utils.Err(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		if httpErr != nil {
			log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}

	responseUserSignUp := utils.Message(true, fmt.Sprintf("User %s created", requestUserSignUpDto.Email))
	log.Info().Msgf("User created: %v", requestUserSignUpDto)
	err = utils.Json(w, http.StatusOK, responseUserSignUp)
	if err != nil {
		log.Error().Msgf("Cannot encode response to json: %s", err)
	}

}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	requestUserSignInDto := dto.RequestUserSignInDto{}

	err := json.NewDecoder(r.Body).Decode(&requestUserSignInDto)
	if err != nil {
		log.Error().Msgf("Error decode json in SignIn() handler: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Error().Msgf("Cannot sending error message to http client from SignIn(): %s", httpErr)
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	responseUserSignIn, err := h.userService.SignIn(ctx, requestUserSignInDto)

	if err != nil {
		if errors.Is(err, userServices.ErrEmailNotValid) {
			log.Error().Msgf("Error user login: %v: %s", requestUserSignInDto, err)
			httpErr := utils.Err(w, http.StatusUnauthorized, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		if errors.Is(err, userServices.ErrEmptyPassword) {
			log.Error().Msgf("Error user login: %v: %s", requestUserSignInDto, err)
			httpErr := utils.Err(w, http.StatusUnauthorized, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		if errors.Is(err, userServices.ErrWrongPassword) {
			log.Error().Msgf("Error user login: %v: %s", requestUserSignInDto, err)
			httpErr := utils.Err(w, http.StatusUnauthorized, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		if errors.Is(err, userServices.ErrUserNotFound) {
			log.Error().Msgf("Error user login: %v: %s", requestUserSignInDto, err)
			httpErr := utils.Err(w, http.StatusUnauthorized, err)
			if httpErr != nil {
				log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
			}
			return
		}

		log.Error().Msgf("Error sign in: %v: %s", requestUserSignInDto, err)
		httpErr := utils.Err(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		if httpErr != nil {
			log.Error().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}

	log.Info().Msgf("User sign in: %v, response: %v", requestUserSignInDto, responseUserSignIn)
	err = utils.Json(w, http.StatusOK, responseUserSignIn)
	if err != nil {
		log.Error().Msgf("Cannot encode response to json: %s", err)
	}

}
