package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
)

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (domain.CalendarUser, error) {
	var userEntity domain.CalendarUser = domain.CalendarUser{}
	tx := r.DB.WithContext(ctx).Limit(1).Where("email = ?", email).Find(&userEntity)
	if tx.Error != nil {
		return domain.CalendarUser{}, fmt.Errorf("error tx in FindUserByEmail(): %w", tx.Error)
	}
	return userEntity, nil
}
func (r *Repository) CreateNewUser(ctx context.Context, user domain.CalendarUser) (domain.CalendarUser, error) {

	findUser := domain.CalendarUser{}
	op := "calendarRepo.CreateNewUser()"

	tx := r.DB.WithContext(ctx).Limit(1).Where("email = ?", user.Email).Find(&findUser)
	if tx.Error != nil {
		return domain.CalendarUser{}, fmt.Errorf("%s: %w", op, tx.Error)
	}

	if (findUser == domain.CalendarUser{}) {
		user.ID = uuid.New()
		tx = r.DB.WithContext(ctx).Create(user)
		if tx.Error != nil {
			return domain.CalendarUser{}, fmt.Errorf("%s: %w", op, tx.Error)
		}
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user domain.CalendarUser) (domain.CalendarUser, error) {

	tx := r.DB.WithContext(ctx).Save(&user)
	if tx.Error != nil {
		return domain.CalendarUser{}, tx.Error
	}
	return user, nil
}
