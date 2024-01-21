package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
)

var (
	errCalendarOrderAlreadyExist error = errors.New("calendar already exist in the database")
)

func (r *Repository) FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error) {
	op := "calendarRepo.FindCalendarByUserId()"
	var cal domain.Calendar

	tx := r.DB.WithContext(ctx).Limit(1).Where("calendar_owner_id = ?", userId).Find(&cal)
	if tx.Error != nil {
		return domain.Calendar{}, fmt.Errorf("%s: %w", op, tx.Error)
	}
	return cal, nil
}

func (r *Repository) CreateCalendar(ctx context.Context, calendar domain.Calendar) (domain.Calendar, error) {
	//findCalendar := domain.Calendar{}
	findCalendar := calendar

	tx := r.DB.WithContext(ctx).Where(findCalendar).FirstOrCreate(&calendar)
	if tx.Error != nil {
		return domain.Calendar{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return domain.Calendar{}, errCalendarOrderAlreadyExist
	}
	return calendar, nil
}

func (r *Repository) UpdateCalendar(ctx context.Context, calendar domain.Calendar) (domain.Calendar, error) {
	tx := r.DB.WithContext(ctx).Save(&calendar)
	if tx.Error != nil {
		return domain.Calendar{}, tx.Error
	}
	return calendar, nil
}
