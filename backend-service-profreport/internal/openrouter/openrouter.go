package openrouter

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"app/main.go/internal/config"
	"app/main.go/internal/models/domain"

	"app/main.go/internal/utils/logger/sl"

	"github.com/google/uuid"
	openrouter "github.com/revrost/go-openrouter"
)

const (
	retryCount    int           = 3
	retryDuration time.Duration = 3 * time.Second
)

type PdfService interface {
	AddJob(
		requestId uuid.UUID,
		inputMarkdown string,
		user domain.User,
	) (chan struct{}, error)
}

type Job struct {
	requestID     uuid.UUID
	questionnaire string
	user          domain.User
	Done          chan struct{}
}

type Openrouter struct {
	logger          *slog.Logger
	cfg             *config.Config
	Client          *openrouter.Client
	pdfService      PdfService
	jobs            chan Job
	shutdownChannel chan struct{}
	wg              *sync.WaitGroup
}

func NewClient(
	logger *slog.Logger,
	cfg *config.Config,
	pdfService PdfService,
) *Openrouter {

	op := "Openrouter.NewClient()"
	log := logger.With(
		slog.String("op", op),
	)

	client := openrouter.NewClient(
		cfg.BotConfig.AI.AIApiToken,
	)

	log.Info("Creating deepseek client")

	return &Openrouter{
		logger:          logger,
		cfg:             cfg,
		Client:          client,
		pdfService:      pdfService,
		jobs:            make(chan Job, cfg.BotConfig.AI.JobBufferSize),
		shutdownChannel: make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}
}

func (s *Openrouter) Start() {
	op := "PdfService.Start()"
	log := s.logger.With(
		slog.String("op", op),
	)
	for i := 0; i < s.cfg.PdfConfig.WorkersCount; i++ {
		s.wg.Add(1)
		go s.handleJob(i)
	}
	log.Info(
		"pdf service started",
	)
	for {
		s.wg.Wait()
	}
}

func (s *Openrouter) AddJob(requestID uuid.UUID, questionnaire string, user domain.User) (chan struct{}, error) {
	newJob := Job{
		requestID:     requestID,
		questionnaire: questionnaire,
		user:          user,
		Done:          make(chan struct{}),
	}
	if len(s.jobs) < s.cfg.PdfConfig.JobBufferSize {
		s.jobs <- newJob
		return newJob.Done, nil
	} else {
		return nil, fmt.Errorf("job buffer is full")
	}
}

func (s *Openrouter) handleJob(id int) {
	defer s.wg.Done()
	op := "Openrouter.handleJob()"
	log := s.logger.With(
		slog.String("op", op),
	)

	log.Info("start openrouter job handler")

	for {

		select {
		case <-s.shutdownChannel:
			return
		case job, ok := <-s.jobs:
			joblog := log.With(
				slog.String("op", op),
				slog.String("requestID", job.requestID.String()),
			)
			if !ok {
				joblog.Error("failed chat completion",
					slog.String("error", "channel is closed"),
				)
				return
			}

			response, err := s.CreateChatCompletion(context.TODO(), joblog, job.requestID, job.questionnaire)

			if err != nil {
				joblog.Error("failed chat completion",
					slog.String("error", err.Error()),
					//slog.String("requestID", job.requestID.String()),
				)
				return
			}

			_, err = s.pdfService.AddJob(job.requestID, response, job.user)
			if err != nil {
				joblog.Error("failed add job to pdf service",
					slog.String("error", err.Error()),
					//slog.String("requestID", job.requestID.String()),
				)
			}

			close(job.Done)

			joblog.Info(
				"AI response received",
				slog.Int("id", id),
				//slog.String("requestID", job.requestID.String()),
			)

		}
	}
}

func (s *Openrouter) CreateChatCompletion(ctx context.Context, logger *slog.Logger, requestId uuid.UUID, message string) (string, error) {
	op := "openrouter.CreateChatCompletion()"
	log := logger.With(
		slog.String("op", op),
		slog.String("requestID", requestId.String()),
	)
	log.Info("start create chat completion")

	//log.Debug("input message", slog.String("message", message))
	var resp openrouter.ChatCompletionResponse
	var err error
	for retry := range retryCount {
		var r openrouter.ChatCompletionResponse
		var e error
		select {
		case <-s.shutdownChannel:
			return "", fmt.Errorf("shutdown openrouter client")
		default:
			r, e = s.Client.CreateChatCompletion(
				ctx,
				openrouter.ChatCompletionRequest{
					Model: s.cfg.BotConfig.AI.ModelName,
					Messages: []openrouter.ChatCompletionMessage{
						openrouter.SystemMessage(s.cfg.BotConfig.AI.SystemRolePromt),
						openrouter.UserMessage(message),
					},
				},
			)
		}
		if e != nil && (isRateLimitError(e) || isEOFError(e)) {
			err = e
			log.Error(
				"Openrouter chat completion response error",
				slog.String("error", err.Error()),
				slog.Int("retry", retry),
			)
			time.Sleep(retryDuration)
			continue
		}
		resp = r
		err = e
		break
	}

	if err != nil {
		// log.Error("error creating chat completion request", slog.String("error", err.Error()))
		return "", fmt.Errorf("error creating chat completion request: %w", err)
	}

	log.Debug("received chat completion response", slog.Any("response role", resp.Choices[0].Message.Role))

	responseText := resp.Choices[0].Message.Content.Text

	err = s.writeResponseInFile(requestId.String(), responseText, "html")
	if err != nil {
		log.Error(
			"error write response in file",
			sl.Err(err),
		)
	}

	return responseText, nil
}

func isRateLimitError(err error) bool {
	// Если библиотека возвращает ошибку с полем Code:
	// if e, ok := err.(interface{ Code() int }); ok && e.Code() == 429 {
	// 	return true
	// }
	// Или проверка по строке ошибки (менее надёжно):
	if err != nil {
		return strings.Contains(err.Error(), "EOF")
	} else {
		return false
	}
}

func isEOFError(err error) bool {
	// Если библиотека возвращает ошибку с полем Code:
	// if e, ok := err.(interface{ Code() int }); ok && e.Code() == 429 {
	// 	return true
	// }
	// Или проверка по строке ошибки (менее надёжно):
	if err != nil {
		return strings.Contains(err.Error(), "HTTP 429")
	} else {
		return false
	}
}

func (s *Openrouter) writeResponseInFile(requestId string, data string, fileType string) error {
	bufWrite := []byte(data)
	filePath := fmt.Sprintf("%s%s.%s", s.cfg.BotConfig.AI.AiResponseFilePath, requestId, fileType)
	err := os.WriteFile(filePath, bufWrite, 0775)
	if err != nil {
		return fmt.Errorf("error write file \"%s\": %w", filePath, err)
	}
	return nil
}

func (s *Openrouter) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("force exit AI client: %w", ctx.Err())
		default:
			close(s.shutdownChannel)
			return nil
		}
	}
}
