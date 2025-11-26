package inMemory

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryCache представляет потокобезопасное хранилище данных в памяти.
// Ключом является ID пользователя, значением - слайс произвольных данных.
type InMemoryCache struct {
	InMemoryMap map[int64][]any
	mutex       sync.RWMutex
}

// NewInMemoryRepository создаёт и возвращает новый экземпляр InMemoryCache.
func NewInMemoryRepository() *InMemoryCache {
	m := make(map[int64][]any)
	return &InMemoryCache{
		InMemoryMap: m,
	}
}

// Save добавляет данные history в кэш для указанного пользователя.
// Если данные для пользователя уже существуют, новые данные добавляются к существующим.
func (r *InMemoryCache) Save(_ context.Context, userID int64, history []any) error {
	if r.InMemoryMap == nil {
		return fmt.Errorf("save: map is not initialized")
	}

	if userID <= 0 {
		return fmt.Errorf("save: invalid userID %d", userID)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.InMemoryMap[userID] = append(r.InMemoryMap[userID], history...)

	return nil
}

// IsUserExist проверяет, существует ли пользователь с указанным ID в кэше.
func (r *InMemoryCache) IsUserExist(_ context.Context, userID int64) (bool, error) {
	if r.InMemoryMap == nil {
		return false, fmt.Errorf("isUserExist: map is not initialized")
	}

	if userID <= 0 {
		return false, fmt.Errorf("isUserExist: invalid userID %d", userID)
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, ok := r.InMemoryMap[userID]
	return ok, nil
}

// Get возвращает данные для указанного пользователя.
// Если пользователь не найден, возвращается ошибка.
func (r *InMemoryCache) Get(_ context.Context, userID int64) ([]any, error) {
	if r.InMemoryMap == nil {
		return nil, fmt.Errorf("get: map is not initialized")
	}

	if userID <= 0 {
		return nil, fmt.Errorf("get: invalid userID %d", userID)
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	val, ok := r.InMemoryMap[userID]
	if !ok {
		return nil, fmt.Errorf("get: userID %d not found", userID)
	}

	return val, nil
}

// Delete удаляет данные пользователя из кэша.
// Если пользователь не существует, операция завершается успешно без ошибки.
func (r *InMemoryCache) Delete(_ context.Context, userID int64) error {
	if r.InMemoryMap == nil {
		return fmt.Errorf("delete: map is not initialized")
	}

	if userID <= 0 {
		return fmt.Errorf("delete: invalid userID %d", userID)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.InMemoryMap, userID)
	return nil
}
