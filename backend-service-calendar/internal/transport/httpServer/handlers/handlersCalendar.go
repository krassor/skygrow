package handlers

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/dto"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/calendar"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils/logger/sl"
	"log/slog"
	"net/http"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type CalendarService interface {
	CreateCalendar(ctx context.Context, calendarUser domain.CalendarUser) (domain.Calendar, error)
	FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error)
}

type CalendarHandler struct {
	calendarService CalendarService
	log             *slog.Logger
}

func NewCalendarHandler(log *slog.Logger, calendarService CalendarService) *CalendarHandler {
	return &CalendarHandler{
		calendarService: calendarService,
		log:             log,
	}
}

func (h *CalendarHandler) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.CreateCalendar()"
	h.log.With(
		slog.String("op", op))

	var httpErr error

	user := r.Context().Value("user").(domain.CalendarUser)
	h.log.Debug("getting user from context", "user", user)

	cal, err := h.calendarService.CreateCalendar(r.Context(), user)
	if err != nil {
		h.log.Error("failed to create calendar", sl.Err(err))

		if errors.Is(err, calendar.ErrCalendarAlreadyExist) {
			httpErr = utils.Err(w, http.StatusConflict, err)
		} else {
			httpErr = utils.Err(w, http.StatusInternalServerError, ErrInternalServer)
		}

		if httpErr != nil {
			h.log.Error("failed to send http answer", sl.Err(err))
		}
		return
	}
	calendarDto := dto.ResponseCalendar{
		CalendarId:       cal.ID.String(),
		CalendarOwnerId:  cal.CalendarOwnerId.String(),
		GoogleCalendarId: cal.GoogleCalendarId,
		Description:      cal.Description,
		Etag:             cal.Etag,
		Summary:          cal.Summary,
		TimeZone:         cal.TimeZone,
		Status:           cal.Status.String(),
	}
	httpErr = utils.Respond(w, utils.Message(true, calendarDto))
	if httpErr != nil {
		h.log.Error("failed to send http answer", sl.Err(err))
	}
	h.log.Debug("calendar created", "calendar", cal)
}

func (h *CalendarHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	op := "httpServer.handlers.CreateCalendar()"
	log := h.log.With(
		slog.String("op", op))

	var httpErr error

	user := r.Context().Value("user").(domain.CalendarUser)
	log.Debug("getting user from context", "user", user)

	cal, err := h.calendarService.FindCalendarByUserId(r.Context(), user.ID)
	if err != nil {
		log.Error("error no calendar founded by user id", sl.Err(err))

		if errors.Is(err, calendar.ErrCalendarNotFound) {
			httpErr = utils.Err(w, http.StatusNotFound, err)
		} else {
			httpErr = utils.Err(w, http.StatusInternalServerError, ErrInternalServer)
		}

		if httpErr != nil {
			log.Error("failed to send http answer", sl.Err(err))
		}
		return
	}
	calendarDto := dto.ResponseCalendar{
		CalendarId:       cal.ID.String(),
		CalendarOwnerId:  cal.CalendarOwnerId.String(),
		GoogleCalendarId: cal.GoogleCalendarId,
		Description:      cal.Description,
		Etag:             cal.Etag,
		Summary:          cal.Summary,
		TimeZone:         cal.TimeZone,
		Status:           cal.Status.String(),
	}
	httpErr = utils.Respond(w, utils.Message(true, calendarDto))
	if httpErr != nil {
		log.Error("failed to send http answer", sl.Err(err))
	}
	log.Debug("calendar founded", "calendar", cal)
}
