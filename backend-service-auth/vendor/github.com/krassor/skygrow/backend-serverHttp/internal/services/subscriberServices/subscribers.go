package subscriberSeervice

import (
	"context"

	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	"github.com/krassor/skygrow/backend-serverHttp/internal/repositories"
)

type SubscriberRepoService interface {
	GetSubscribers(ctx context.Context) ([]entities.Subscriber, error)
	CreateNewSubscriber(ctx context.Context, chatId int64, name string) (entities.Subscriber, error)
	GetSubscriberByChatId(ctx context.Context, chatId int64) (entities.Subscriber, error)
	UpdateSubscriberByChatId(ctx context.Context, subscriber entities.Subscriber, subscriberStatus bool) (entities.Subscriber, error)
}

type subscriberRepoService struct {
	subscriberRepository repositories.SubscribersRepo
}

func NewSubscriberRepoService(subscriberRepository repositories.SubscribersRepo) SubscriberRepoService {
	return &subscriberRepoService{
		subscriberRepository: subscriberRepository,
	}
}

func (s *subscriberRepoService) GetSubscribers(ctx context.Context) ([]entities.Subscriber, error) {
	subscriber, err := s.subscriberRepository.FindAllSubscribers(ctx)
	return subscriber, err
}

func (s *subscriberRepoService) GetSubscriberByChatId(ctx context.Context, chatId int64) (entities.Subscriber, error) {
	subscriber, err := s.subscriberRepository.FindSubscriberByChatId(ctx, chatId)
	return subscriber, err
}

func (s *subscriberRepoService) CreateNewSubscriber(ctx context.Context, chatId int64, name string) (entities.Subscriber, error) {

	subscriber, err := s.subscriberRepository.CreateSubscriber(ctx, chatId, name)
	return subscriber, err
}

func (s *subscriberRepoService) UpdateSubscriberByChatId(ctx context.Context, subscriber entities.Subscriber, subscriberStatus bool) (entities.Subscriber, error) {
	subscriber.IsActive = subscriberStatus
	subscriberResponse, err := s.subscriberRepository.UpdateSubscriber(ctx, subscriber)
	return subscriberResponse, err
}

