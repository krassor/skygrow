package openrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"app/main.go/internal/config"
	"app/main.go/internal/models/domain"
	"app/main.go/internal/models/dto"

	"app/main.go/internal/utils/logger/sl"

	"github.com/google/uuid"
	openrouter "github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

const (
	// retryCount определяет количество попыток повторного запроса при ошибках.
	retryCount int = 3
	// retryDuration задаёт интервал между попытками повторного запроса.
	retryDuration time.Duration = 3 * time.Second
)

// PdfService — интерфейс, который должен реализовывать сервис для генерации PDF.
// Используется для инъекции зависимости в Openrouter.
type PdfService interface {
	// AddJob добавляет задачу на генерацию PDF.
	// Принимает UUID запроса, строку с контентом и данные пользователя.
	// Возвращает канал для сигнализации завершения и ошибку (если есть).
	AddJob(
		requestId uuid.UUID,
		input string,
		user domain.User,
		jobType string,
	) (chan struct{}, error)
}

// Job представляет задачу, передаваемую в воркер.
type Job struct {
	requestID     uuid.UUID     // Уникальный идентификатор запроса
	jobType       string        // Тип запроса
	questionnaire string        // Текст анкеты для обработки ИИ
	user          domain.User   // Данные пользователя
	Done          chan struct{} // Канал для сигнала завершения
}

// Openrouter — структура, управляющая взаимодействием с OpenRouter API.
// Содержит пул воркеров для асинхронной обработки запросов.
type Openrouter struct {
	logger          *slog.Logger       // Логгер с контекстом
	cfg             *config.Config     // Конфигурация приложения
	Client          *openrouter.Client // Клиент OpenRouter API
	pdfService      PdfService         // Сервис для генерации PDF
	jobs            chan Job           // Канал задач
	shutdownChannel chan struct{}      // Канал для сигнала завершения
	wg              *sync.WaitGroup    // Группа для ожидания завершения воркеров
}

// NewClient создаёт новый экземпляр Openrouter.
//
// Параметры:
//   - logger: экземпляр *slog.Logger для логирования.
//   - cfg: конфигурация приложения.
//   - pdfService: реализация PdfService для генерации PDF.
//
// Возвращает указатель на инициализированный Openrouter.
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

	log.Info("Creating openrouter client")

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

// Start запускает воркеры для обработки задач.
// Количество воркеров задаётся в конфиге (WorkersCount).
// Метод блокируется до завершения всех воркеров.
func (s *Openrouter) Start() {
	op := "Openrouter.Start()"
	log := s.logger.With(
		slog.String("op", op),
	)
	for i := 0; i < s.cfg.BotConfig.AI.WorkersCount; i++ {
		s.wg.Add(1)
		go s.handleJob(i)
	}
	log.Info("openrouter service started")

	s.wg.Wait()
}

// AddJob добавляет новую задачу в очередь на обработку.
//
// Параметры:
//   - requestID: уникальный идентификатор запроса.
//   - questionnaire: текст анкеты для анализа.
//   - user: данные пользователя.
//   - jobType: тип запроса: ADULT, SCHOOLCHILD
//
// Возвращает:
//   - Канал Done для ожидания завершения задачи.
//   - Ошибку, если буфер задач переполнен.
//
// Если буфер заполнен, задача не добавляется.
func (s *Openrouter) AddJob(requestID uuid.UUID, questionnaire string, user domain.User, jobType string) (chan struct{}, error) {
	newJob := Job{
		requestID:     requestID,
		questionnaire: questionnaire,
		user:          user,
		jobType:       jobType,
		Done:          make(chan struct{}),
	}
	select {
	case <-s.shutdownChannel:
		return nil, fmt.Errorf("service is shutting down")
	default:

		if len(s.jobs) < s.cfg.BotConfig.AI.JobBufferSize {
			s.jobs <- newJob
			return newJob.Done, nil
		} else {
			return nil, fmt.Errorf("job buffer is full")
		}
	}
}

// handleJob — воркер, обрабатывающий задачи из канала.
// Получает задачу, отправляет запрос в OpenRouter, сохраняет ответ и передаёт в PDF-сервис.
// Цикл прерывается при закрытии канала заданий или сигнале завершения.
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
			defer close(job.Done)
			requestID := job.requestID
			user := job.user
			questionnaire := job.questionnaire

			joblog := log.With(
				slog.String("op", op),
				slog.String("requestID", requestID.String()),
			)
			if !ok {
				joblog.Error("failed chat completion",
					slog.String("error", "channel is closed"),
				)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), s.cfg.BotConfig.AI.GetTimeout())

			response, err := s.CreateChatCompletion(ctx, joblog, requestID, questionnaire, job.jobType)

			cancel()
			if err != nil {
				joblog.Error("failed chat completion",
					slog.String("error", err.Error()),
				)
				continue
			}

			_, err = s.pdfService.AddJob(requestID, response, user, job.jobType)
			if err != nil {
				joblog.Error("failed add job to pdf service",
					slog.String("error", err.Error()),
				)
				continue
			}

			joblog.Info(
				"AI response received",
				slog.Int("id", id),
			)
		}
	}
}

// CreateChatCompletion отправляет запрос в OpenRouter для генерации ответа.
//
// Параметры:
//   - ctx: контекст с таймаутом.
//   - logger: логгер с контекстом.
//   - requestId: UUID запроса.
//   - message: пользовательское сообщение (анкета).
//   - jobType: тип запроса: ADULT, SCHOOLCHILD
//
// Возвращает:
//   - Строка с ответом ИИ.
//   - Ошибка, если запрос не удался.
//
// Повторяет запрос до retryCount раз при ошибках 429 или EOF.
func (s *Openrouter) CreateChatCompletion(ctx context.Context, logger *slog.Logger, requestId uuid.UUID, message string, jobType string) (string, error) {
	op := "openrouter.CreateChatCompletion()"
	log := logger.With(
		slog.String("op", op),
		slog.String("requestID", requestId.String()),
	)
	log.Info("start create chat completion")

	var resp openrouter.ChatCompletionResponse
	var err error
	var prompt string

	switch jobType {
	case "ADULT":
		prompt = s.cfg.BotConfig.AI.AdultSystemRolePrompt
	case "SCHOOLCHILD":
		prompt = s.cfg.BotConfig.AI.SchoolchildSystemRolePrompt
	default:
		log.Error(
			"unknown job type",
		)
		return "", fmt.Errorf("unknown job type")
	}

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
						openrouter.SystemMessage(prompt),
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
		return "", fmt.Errorf("return error creating chat completion request: %w", err)
	}

	log.Debug("received chat completion response", slog.Any("response role", resp.Choices[0].Message.Role))

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty AI response (no resp.Choices)")
	}
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

func (s *Openrouter) CreateChatCompletionWithStructuredOutputs(ctx context.Context, logger *slog.Logger, requestId uuid.UUID, message string, jobType string) (string, error) {
	op := "openrouter.CreateChatCompletion()"
	log := logger.With(
		slog.String("op", op),
		slog.String("requestID", requestId.String()),
	)
	log.Info("start create chat completion")

	//log.Debug("input message", slog.String("message", message))
	var responseSchema dto.StructuredResponseSchema
	var resp openrouter.ChatCompletionResponse
	var err error
	var prompt string

	switch jobType {
	case "ADULT":
		prompt = s.cfg.BotConfig.AI.StructuredOutputs.AdultSOSystemRolePrompt
	case "SCHOOLCHILD":
		prompt = s.cfg.BotConfig.AI.StructuredOutputs.SchoolchildSOSystemRolePrompt
	default:
		log.Error(
			"unknown job type",
		)
		return "", fmt.Errorf("unknown job type")
	}

	for retry := range retryCount {
		var r openrouter.ChatCompletionResponse
		var e error
		select {
		case <-s.shutdownChannel:
			return "", fmt.Errorf("shutdown openrouter client")
		default:

			schema, err := jsonschema.GenerateSchemaForType(responseSchema)
			if err != nil {
				log.Error(
					"GenerateSchemaForType error",
					sl.Err(err),
				)
				return "", fmt.Errorf("GenerateSchemaForType error: %w", err)
			}

			r, e = s.Client.CreateChatCompletion(
				ctx,
				openrouter.ChatCompletionRequest{
					Model: s.cfg.BotConfig.AI.ModelName,
					Messages: []openrouter.ChatCompletionMessage{
						openrouter.SystemMessage(prompt),
						openrouter.UserMessage(message),
					},
					ResponseFormat: &openrouter.ChatCompletionResponseFormat{
						Type: "json_schema",
						JSONSchema: &openrouter.ChatCompletionResponseFormatJSONSchema{
							Name:   "nameTODO",
							Strict: true,
							Schema: schema,
						},
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
		return "", fmt.Errorf("return error creating chat completion request: %w", err)
	}

	log.Debug("received chat completion response", slog.Any("response: ", resp))

	//responseText := resp.Choices[0].Message.Content.Text
	b := []byte(resp.Choices[0].Message.Content.Text)
	err = json.Unmarshal(b, &responseSchema)
	if err != nil {
		log.Error(
			"error unmarshal response",
			sl.Err(err),
		)
	}
	log.Debug(
		"response schema",
		slog.Any("schema", responseSchema),
	)
	responseText := responseSchema.InterpretationRIASEC + "/r/n" + responseSchema.InterpretationRIASEC

	err = s.writeResponseInFile(requestId.String(), responseText, "html")
	if err != nil {
		log.Error(
			"error write response in file",
			sl.Err(err),
		)
	}

	return responseText, nil
}

// isRateLimitError проверяет, связана ли ошибка с превышением лимита запросов (HTTP 429).
// Временное решение по анализу строки ошибки — менее надёжно, чем проверка кода.
func isRateLimitError(err error) bool {
	if err != nil {
		//return strings.Contains(err.Error(), "HTTP 429")
		return strings.Contains(err.Error(), "429")
	}
	return false
}

// isEOFError проверяет, связана ли ошибка с разрывом соединения (EOF).
// Используется для повтора запроса.
func isEOFError(err error) bool {
	if err != nil {
		return strings.Contains(err.Error(), "EOF")
	}
	return false
}

// writeResponseInFile сохраняет текстовый ответ ИИ в файл.
//
// Параметры:
//   - requestId: идентификатор запроса (часть имени файла).
//   - data: содержимое для записи.
//   - fileType: расширение файла (например, "html").
//
// Использует filepath.Clean для защиты от path traversal.
// Устанавливает права 0644.
// Возвращает ошибку при неудачной записи.
func (s *Openrouter) writeResponseInFile(requestId string, data string, fileType string) error {
	if _, err := uuid.Parse(requestId); err != nil {
		return fmt.Errorf("invalid requestId")
	}
	bufWrite := []byte(data)
	filePath := filepath.Clean(fmt.Sprintf("%s%s.%s", s.cfg.BotConfig.AI.AiResponseFilePath, requestId, fileType))
	err := os.WriteFile(filePath, bufWrite, 0644)
	if err != nil {
		return fmt.Errorf("error write file \"%s\": %w", filePath, err)
	}
	return nil
}

// Shutdown корректно завершает работу сервиса.
//
// Параметры:
//   - ctx: контекст для отслеживания таймаута завершения.
//
// Действия:
//   - Закрывает канал shutdownChannel.
//   - Закрывает канал jobs для остановки воркеров.
//   - Возвращает ошибку, если контекст отменён.
//
// После вызова новые задачи не принимаются, обработка текущих завершается.
func (s *Openrouter) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("force exit AI client: %w", ctx.Err())
	default:
		close(s.shutdownChannel)
		close(s.jobs)
		return nil
	}
}
