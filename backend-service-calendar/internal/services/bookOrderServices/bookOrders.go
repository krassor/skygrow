package bookOrderServices

// import (
// 	"context"
// 	"time"

// 	"github.com/google/uuid"

// 	"github.com/krassor/skygrow/backend-service-calendar/internal/models/dto"
// 	"github.com/krassor/skygrow/backend-service-calendar/internal/models/entities"
// 	telegramBot "github.com/krassor/skygrow/backend-service-calendar/internal/telegram"
// )

// type BookOrderRepository interface {
// 	FindAllBookOrder(ctx context.Context) ([]entities.BookOrder, error)
// 	CreateBookOrder(ctx context.Context, bookOrder entities.BookOrder) (entities.BookOrder, error)
// 	UpdateBookOrder(ctx context.Context, bookOrder entities.BookOrder) (entities.BookOrder, error)
// 	FindBookOrderById(ctx context.Context, id string) (entities.BookOrder, error)
// }

// type BookOrderService struct {
// 	BookOrderRepository BookOrderRepository
// 	tgBot               *telegramBot.Bot
// }

// func NewBookOrderService(r BookOrderRepository, tgBot *telegramBot.Bot) *BookOrderService {
// 	return &BookOrderService{
// 		BookOrderRepository: r,
// 		tgBot:               tgBot,
// 	}
// }

// func (s *BookOrderService) CreateNewBookOrder(ctx context.Context, bookOrderDto dto.RequestBookOrderDto) (string, error) {

// 	bookOrderEntity := entities.BookOrder{
// 		FirstName:          bookOrderDto.FirstName,
// 		SecondName:         bookOrderDto.SecondName,
// 		Phone:              bookOrderDto.Phone,
// 		Email:              bookOrderDto.Email,
// 		MentorID:           bookOrderDto.MentorID,
// 		ProblemDescription: bookOrderDto.ProblemDescription,
// 	}
// 	bookOrderEntity.BaseModel.ID = uuid.NewString()
// 	bookOrderEntity.BaseModel.Created_at = time.Now()
// 	bookOrderEntity.BaseModel.Updated_at = bookOrderEntity.BaseModel.Created_at

// 	responseBookOrderEntity, err := s.BookOrderRepository.CreateBookOrder(ctx, bookOrderEntity)
// 	if err != nil {
// 		return "", err
// 	}

// 	ResponseBookOrderDto := dto.ResponseBookOrderDto{
// 		FirstName:          responseBookOrderEntity.FirstName,
// 		SecondName:         responseBookOrderEntity.SecondName,
// 		Phone:              responseBookOrderEntity.Phone,
// 		Email:              responseBookOrderEntity.Email,
// 		MentorID:           responseBookOrderEntity.MentorID,
// 		BookOrderID:        responseBookOrderEntity.BaseModel.ID,
// 		ProblemDescription: responseBookOrderEntity.ProblemDescription,
// 	}
// 	err = s.tgBot.BookOrderNotify(ctx, ResponseBookOrderDto)

// 	if err != nil {
// 		return "", err
// 	}

// 	return responseBookOrderEntity.BaseModel.ID, nil
// }
