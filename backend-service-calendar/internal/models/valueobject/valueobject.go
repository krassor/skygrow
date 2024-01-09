package valueobject

import (
	"github.com/google/uuid"
	"time"
)

type SubscribeCalendarEvent struct {
	ID         uuid.UUID
	CalendarId uuid.UUID
	AttendeeId uuid.UUID
	CreatedAt  time.Time
}

type ChangeCalendar struct {
	ID         uuid.UUID
	CalendarId uuid.UUID
	UserId     uuid.UUID
	CreatedAt  time.Time
}
