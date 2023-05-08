package inMemory

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type InMemoryRepository struct {
	InMemoryMap map[string][]openai.ChatCompletionMessage
}

func NewInMemoryRepository() *InMemoryRepository {
	m := make(map[string][]openai.ChatCompletionMessage)
	return &InMemoryRepository{
		InMemoryMap: m,
	}
}

func (r *InMemoryRepository) SaveUserMessage(ctx context.Context, username string, message openai.ChatCompletionMessage) error {
	if r.InMemoryMap == nil {
		return fmt.Errorf("SaveUserMessage error: Map is not initializate")
	}

	if username == "" {
		return fmt.Errorf("SaveUserMessage error: Empty key \"username\"")
	}

	r.InMemoryMap[username] = append(r.InMemoryMap[username], message)

	return nil
}

func (r *InMemoryRepository) IsUserExist(ctx context.Context, username string) (bool, error) {
	if r.InMemoryMap == nil {
		return false, fmt.Errorf("IsUserExist error: Map is not initializate")
	}

	if username == "" {
		return false, fmt.Errorf("IsUserExist error: Empty key \"username\"")
	}

	_, ok := r.InMemoryMap[username]
	if ok {
		return true, nil
	} else {
		return false, nil
	}

}

func (r *InMemoryRepository) LoadUserMessages(ctx context.Context, username string) ([]openai.ChatCompletionMessage, error) {
	if r.InMemoryMap == nil {
		return nil, fmt.Errorf("LoadUserMessages error: Map is not initializate")
	}

	if username == "" {
		return nil, fmt.Errorf("LoadUserMessages error: Empty key \"username\"")
	}

	val, ok := r.InMemoryMap[username]
	if ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("LoadUserMessages error: Username not found")
	}
}
