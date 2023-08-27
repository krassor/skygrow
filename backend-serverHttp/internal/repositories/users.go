package repositories

import (
	"context"
	"fmt"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
)

func (r *repository) FindUserByEmail(ctx context.Context, email string) (entities.User, error) {
	var userEntity entities.User = entities.User{}
	tx := r.DB.WithContext(ctx).Limit(1).Where("email = ?", email).Find(&userEntity)
	if tx.Error != nil {
		return entities.User{}, fmt.Errorf("error tx in FindUserByEmail(): %w", tx.Error)
	}
	return userEntity, nil
}
func (r *repository) CreateNewUser(ctx context.Context, user entities.User) (entities.User, error) {

	tx := r.DB.WithContext(ctx).Create(&user)
	if tx.Error != nil {
		return entities.User{}, fmt.Errorf("error tx in CreateNewUser(): %w", tx.Error)
	}
	return user, nil
}
