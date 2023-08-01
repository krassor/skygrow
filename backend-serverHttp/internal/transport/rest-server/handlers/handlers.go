package handlers

import (
	"encoding/json"

	"net/http"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/dto"
	"github.com/krassor/skygrow/backend-serverHttp/internal/services/bookOrderServices"
	"github.com/krassor/skygrow/backend-serverHttp/pkg/utils"
	"github.com/rs/zerolog/log"
)

type BookOrderHandler struct {
	bookOrderService *bookOrderServices.BookOrderService
}

func NewBookOrderHandler(bookOrderService *bookOrderServices.BookOrderService) *BookOrderHandler {
	return &BookOrderHandler{
		bookOrderService: bookOrderService,
	}
}

func (h *BookOrderHandler) CreateBookOrder(w http.ResponseWriter, r *http.Request) {
	bookOrderDto := dto.RequestBookOrderDto{}

	err := json.NewDecoder(r.Body).Decode(&bookOrderDto)
	if err != nil {
		log.Warn().Msgf("Error decode json: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Warn().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}

	bookOrderID, err := h.bookOrderService.CreateNewBookOrder(r.Context(), bookOrderDto)

	if err != nil {
		log.Error().Msgf("Error creating book order: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Warn().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}

	responseBookOrderDto := dto.ResponseBookOrderDto{
		FirstName:   bookOrderDto.FirstName,
		SecondName:  bookOrderDto.SecondName,
		Phone:       bookOrderDto.Phone,
		Email:       bookOrderDto.Email,
		MentorID:    bookOrderDto.MentorID,
		BookOrderID: bookOrderID,
	}

	responseBookOrder := utils.Message(true, responseBookOrderDto)
	err = utils.Json(w, http.StatusOK, responseBookOrder)
	if err != nil {
		log.Warn().Msgf("Cannot encode response to json: %s", err)
	}
}
