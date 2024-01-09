package calendar

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/GoogleService"
)

var (
	ErrEmailNotValid    = errors.New("email is not valid")
	ErrNilEntity        = errors.New("entity is nil")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrCalendarNotFound = errors.New("calendar not found")
	ErrWrongPassword    = errors.New("wrong password")
)

type calendarRepository interface {
	FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error)
	CreateCalendar(ctx context.Context, calendar domain.Calendar) (domain.Calendar, error)
	UpdateCalendar(ctx context.Context, calendar domain.Calendar) (domain.Calendar, error)
}

type Calendar struct {
	calendarRepository calendarRepository
	googleCalendar     *GoogleService.GoogleCalendar
}

func NewCalendarService(cr calendarRepository, gc *GoogleService.GoogleCalendar) *Calendar {
	return &Calendar{
		calendarRepository: cr,
		googleCalendar:     gc,
	}
}

func (c *Calendar) FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error) {
	op := "Calendar FindCalendarByUserId()"
	calendar, err := c.calendarRepository.FindCalendarByUserId(ctx, userId)
	if err != nil {
		return domain.Calendar{}, fmt.Errorf("%s:%w", op, err)
	}
	if (calendar == domain.Calendar{}) {
		return domain.Calendar{}, fmt.Errorf("%s:%w", op, ErrCalendarNotFound)
	}
	return calendar, nil
}

func (c *Calendar) CreateCalendar(ctx context.Context, calendar domain.Calendar) (domain.Calendar, error) {
	op := "Calendar CreateCalendar()"
	if (calendar == domain.Calendar{}) {
		return domain.Calendar{}, fmt.Errorf("%s:%w", op, ErrNilEntity)
	}

	googleCalendarId, err := c.googleCalendar.CreateCalendar(calendar.Description, calendar.Summary, calendar.TimeZone)
	if err != nil {
		return domain.Calendar{}, fmt.Errorf("%s : %w", op, err)
	}

	calendarOut, err := c.calendarRepository.CreateCalendar(ctx, calendar)
	if err != nil {
		return domain.Calendar{}, fmt.Errorf("%s:%w", op, err)
	}

	return calendarOut, nil
}
