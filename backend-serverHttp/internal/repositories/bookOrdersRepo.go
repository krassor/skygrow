package repositories

import (
	"context"
	"errors"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
)

var (
	errBookOrderAlreadyExist error = errors.New("bookOrder already exist in the database")
)

func (r *repository) FindAllBookOrder(ctx context.Context) ([]entities.BookOrder, error) {
	var bookOrders []entities.BookOrder
	tx := r.DB.WithContext(ctx).Find(&bookOrders)
	if tx.Error != nil {
		return []entities.BookOrder{}, tx.Error
	}

	return bookOrders, nil
}

func (r *repository) FindBookOrderById(ctx context.Context, id string) (entities.BookOrder, error) {
	var bookOrder entities.BookOrder
	tx := r.DB.WithContext(ctx).First(&bookOrder, id)
	if tx.Error != nil {
		return entities.BookOrder{}, tx.Error
	}
	return bookOrder, nil
}

func (r *repository) CreateBookOrder(ctx context.Context, bookOrder entities.BookOrder) (entities.BookOrder, error) {

	findBookOrder := entities.BookOrder{}
	findBookOrder.BaseModel.ID = bookOrder.BaseModel.ID

	tx := r.DB.WithContext(ctx).Where(findBookOrder).FirstOrCreate(&bookOrder)
	if tx.Error != nil {
		return entities.BookOrder{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entities.BookOrder{}, errBookOrderAlreadyExist
	}
	return bookOrder, nil
}

func (r *repository) UpdateBookOrder(ctx context.Context, bookOrder entities.BookOrder) (entities.BookOrder, error) {
	tx := r.DB.WithContext(ctx).Save(&bookOrder)
	if tx.Error != nil {
		return entities.BookOrder{}, tx.Error
	}
	return bookOrder, nil
}
