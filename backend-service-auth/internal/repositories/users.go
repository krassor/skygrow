package repositories

import (
	"context"
	"fmt"

	"github.com/krassor/skygrow/backend-service-auth/internal/models/domain"
)

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var userEntity domain.User = domain.User{}
	tx := r.DB.WithContext(ctx).Limit(1).Where("email = ?", email).Find(&userEntity)
	if tx.Error != nil {
		return domain.User{}, fmt.Errorf("error tx in FindUserByEmail(): %w", tx.Error)
	}
	return userEntity, nil
}
func (r *Repository) CreateNewUser(ctx context.Context, user domain.User) (domain.User, error) {

	tx := r.DB.WithContext(ctx).Create(&user)
	if tx.Error != nil {
		return domain.User{}, fmt.Errorf("error tx in CreateNewUser(): %w", tx.Error)
	}
	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {

	tx := r.DB.WithContext(ctx).Save(&user)
	if tx.Error != nil {
		return domain.User{}, tx.Error
	}
	return user, nil
}
