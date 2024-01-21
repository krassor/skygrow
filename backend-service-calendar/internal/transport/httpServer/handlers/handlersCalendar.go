package handlers

import (
	"context"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils/logger/sl"
	"log/slog"
	"net/http"
)

type CalendarService interface {
	CreateCalendar(ctx context.Context, calendarUser domain.User) (domain.Calendar, error)
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

	user := r.Context().Value("user").(domain.User)

	cal, err := h.calendarService.CreateCalendar(r.Context(), user)
	if err != nil {
		h.log.Error("failed to create calendar", sl.Err(err))
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			h.log.Error("failed to send http answer", sl.Err(err))
		}
		return
	}

	h.log.Debug("calendar created", "calendar", cal)
}
