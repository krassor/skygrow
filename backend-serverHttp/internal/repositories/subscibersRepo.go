package repositories

import (
	"context"
	"errors"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
)

var (
	errSubscriberAlreadyExist error = errors.New("subscriber already exist in the database")
)

type SubscribersRepo interface {
	FindAllSubscribers(ctx context.Context) ([]entities.Subscriber, error)
	CreateSubscriber(ctx context.Context, chatId int64, name string) (entities.Subscriber, error)
	UpdateSubscriber(ctx context.Context, subscriber entities.Subscriber) (entities.Subscriber, error)
	FindSubscriberByChatId(ctx context.Context, chatId int64) (entities.Subscriber, error)
}

func (r *repository) FindAllSubscribers(ctx context.Context) ([]entities.Subscriber, error) {
	var subscribers []entities.Subscriber
	tx := r.DB.WithContext(ctx).Find(&subscribers)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return subscribers, nil
}

func (r *repository) CreateSubscriber(ctx context.Context, chatId int64, name string) (entities.Subscriber, error) {
	var subscriber entities.Subscriber = entities.Subscriber{
		Name:     name,
		ChatID:   chatId,
		IsActive: true,
	}
	tx := r.DB.WithContext(ctx).Where(entities.Subscriber{ChatID: chatId}).FirstOrCreate(&subscriber)
	if tx.Error != nil {
		return entities.Subscriber{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entities.Subscriber{}, errSubscriberAlreadyExist
	}
	return subscriber, nil
}

func (r *repository) UpdateSubscriber(ctx context.Context, subscriber entities.Subscriber) (entities.Subscriber, error) {

	tx := r.DB.WithContext(ctx).Save(&subscriber)
	if tx.Error != nil {
		return entities.Subscriber{}, tx.Error
	}
	return subscriber, nil
}

func (r *repository) FindSubscriberByChatId(ctx context.Context, chatId int64) (entities.Subscriber, error) {
	var subscriber entities.Subscriber

	tx := r.DB.WithContext(ctx).Limit(1).Where("chatId = ?", chatId).Find(&subscriber)
	if tx.Error != nil {
		return entities.Subscriber{}, tx.Error
	}
	return subscriber, nil
}
