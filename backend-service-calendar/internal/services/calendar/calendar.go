package calendar

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/GoogleService"
	"log/slog"
)

type calendarRepository interface {
	FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error)
	CreateCalendar(ctx context.Context, calendar *domain.Calendar) error
	UpdateCalendar(ctx context.Context, calendar *domain.Calendar) error
}

type Calendar struct {
	calendarRepository calendarRepository
	googleCalendar     *GoogleService.GoogleCalendar
	log                *slog.Logger
}

func NewCalendarService(log *slog.Logger, cr calendarRepository, gc *GoogleService.GoogleCalendar) *Calendar {
	return &Calendar{
		calendarRepository: cr,
		googleCalendar:     gc,
		log:                log,
	}
}

func (c *Calendar) FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error) {
	op := "Calendar FindCalendarByUserId()"
	log := c.log.With(
		slog.String("op", op))

	calendar, err := c.calendarRepository.FindCalendarByUserId(ctx, userId)

	if errors.Is(err, ErrCalendarNotFound) {
		return domain.Calendar{}, err
	}

	if err != nil {
		log.Error("calendar module error", "err", err)
		return domain.Calendar{}, err
	}

	if (calendar == domain.Calendar{}) {
		log.Warn("calendar not found, but no error from repository module")
		return domain.Calendar{}, ErrCalendarNotFound
	}
	return calendar, nil
}

func (c *Calendar) CreateCalendar(ctx context.Context, calendarUser domain.CalendarUser) (domain.Calendar, error) {
	op := "services.calendar.CreateCalendar()"
	log := c.log.With(
		slog.String("op", op))

	log.Debug("input user", "user", calendarUser)
	if (calendarUser == domain.CalendarUser{}) {
		return domain.Calendar{}, fmt.Errorf("%s:%w", op, ErrNilEntity)
	}

	summary := fmt.Sprintf("%s %s", calendarUser.FirstName, calendarUser.SecondName)
	description := fmt.Sprintf("calendar owner id: %s", calendarUser.ID.String())

	//!!!!!!!!!!!!!!
	timeZone := "Europe/Moscow" //TODO: timeZone should be received from frontend

	findCal, err := c.FindCalendarByUserId(ctx, calendarUser.ID)
	if (err == nil) && (findCal != domain.Calendar{}) {
		log.Debug("calendar already exist", "calendar", findCal)
		return domain.Calendar{}, ErrCalendarAlreadyExist
	}

	googleCalendarId, err := c.googleCalendar.CreateCalendar(description, summary, timeZone)
	if err != nil {
		return domain.Calendar{}, fmt.Errorf("%s : %w", op, err)
	}

	calendar := domain.Calendar{
		CalendarOwnerId:  calendarUser.ID,
		GoogleCalendarId: googleCalendarId,
		Description:      description,
		Summary:          summary,
		TimeZone:         timeZone,
	}
	log.Debug("calendar with Google ID", "calendar", calendar)

	err = c.calendarRepository.CreateCalendar(ctx, &calendar)
	if err != nil {
		log.Debug("error repo createCalendar", "err", err)
		return domain.Calendar{}, fmt.Errorf("%s:%w", op, err)
	}
	log.Debug("calendarOut", "calendar", calendar)
	return calendar, nil
}
