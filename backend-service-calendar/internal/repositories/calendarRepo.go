package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/calendar"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils/logger/sl"
	"gorm.io/gorm"
	"log/slog"
)

var (
	ErrCalendarAlreadyExist error = errors.New("calendar already exist in the database")
)

func (r *Repository) FindCalendarByUserId(ctx context.Context, userId uuid.UUID) (domain.Calendar, error) {
	op := "calendarRepo.FindCalendarByUserId()"
	log := r.log.With(
		slog.String("op", op))
	var cal domain.Calendar

	tx := r.DB.WithContext(ctx).Limit(1).Where("calendar_owner_id = ?", userId).Find(&cal)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		log.Debug("", sl.Err(tx.Error))
		return domain.Calendar{}, calendar.ErrCalendarNotFound
	}
	if tx.Error != nil {
		log.Debug("", sl.Err(tx.Error))
		return domain.Calendar{}, fmt.Errorf("%s: %w", op, tx.Error)
	}
	return cal, nil
}

func (r *Repository) CreateCalendar(ctx context.Context, cal *domain.Calendar) error {
	findCalendar := domain.Calendar{}
	op := "calendarRepo.CreateCalendar()"
	log := r.log.With(
		slog.String("op", op))

	log.Debug("input calendar", "calendar", *cal)
	tx := r.DB.WithContext(ctx).Limit(1).Where("calendar_owner_id = ?", cal.CalendarOwnerId).Find(&findCalendar)
	if tx.Error != nil {
		log.Debug("tx error", "err", tx.Error)
		return fmt.Errorf("%s: %w", op, tx.Error)
	}
	log.Debug("find cal in db", "calendar", findCalendar)

	if (findCalendar == domain.Calendar{}) {
		log.Debug("findCalendar = domain.Calendar{}")
		cal.ID = uuid.New()
		log.Debug("create calendar", "calendar", cal)
		tx := r.DB.WithContext(ctx).Create(cal)
		if tx.Error != nil {
			return fmt.Errorf("%s: %w", op, tx.Error)
		}
		return nil
	}

	return ErrCalendarAlreadyExist
}

func (r *Repository) UpdateCalendar(ctx context.Context, cal *domain.Calendar) error {
	op := "calendarRepo.UpdateCalendar()"
	log := r.log.With(
		slog.String("op", op))

	tx := r.DB.WithContext(ctx).Save(&cal)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		log.Debug("", sl.Err(tx.Error))
		return calendar.ErrCalendarNotFound
	}

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
