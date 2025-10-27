package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"app/main.go/internal/config"
	"app/main.go/internal/models/domain"

	"github.com/google/uuid"
	"github.com/nativebpm/gotenberg"
)

const (
	retryCount    int           = 1
	retryDuration time.Duration = 3 * time.Second
	timeout       time.Duration = 30 * time.Second
)

type ViewData struct {
	User           domain.User
	AiHtmlResponse string
}

type MailService interface {
	AddJob(
		requestID uuid.UUID,
		user domain.User,
		subject string,
	) error
}

type Job struct {
	requestID uuid.UUID
	input     string
	user      domain.User
	Done      chan struct{}
}

// Mailer представляет собой клиент для отправки электронных писем.
type PdfService struct {
	logger          *slog.Logger
	cfg             *config.Config
	jobs            chan Job
	mailService     MailService
	shutdownChannel chan struct{}
	wg              *sync.WaitGroup
}

// NewMailer создает новый экземпляр Mailer.
// Проверяет, что все необходимые параметры присутствуют в конфиге.
func New(
	logger *slog.Logger,
	cfg *config.Config,
	mailService MailService,
) *PdfService {
	return &PdfService{
		logger:          logger,
		cfg:             cfg,
		mailService:     mailService,
		jobs:            make(chan Job, cfg.PdfConfig.JobBufferSize),
		shutdownChannel: make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}
}

func (m *PdfService) Start() {
	op := "PdfService.Start()"
	log := m.logger.With(
		slog.String("op", op),
	)
	for i := 0; i < m.cfg.PdfConfig.WorkersCount; i++ {
		m.wg.Add(1)
		go m.handleJob(i)
	}
	log.Info(
		"pdf service started",
	)
	for {
		m.wg.Wait()
	}
}

// AddJob добавляет новую задачу на конвертацию Markdown в PDF в очередь.
//
// Параметры:
//   - requestID: уникальный идентификатор запроса (UUID).
//   - inputMarkdown: строка с содержимым Markdown, которое необходимо обработать.
//
// Возвращает:
//   - Канал `chan struct{}` для отслеживания завершения обработки задачи.
//     Клиент должен слушать этот канал, чтобы получить сигнал о завершении.
//   - Ошибку `error`, если буфер задач переполнен (возвращается `nil` для канала).
//
// Логика работы:
//  1. Создаётся новая задача (`Job`) с указанным `requestID` и `inputMarkdown`.
//  2. Если в буфере (`m.jobs`) есть свободное место, задача отправляется в канал,
//     а клиент получает канал `Done` для отслеживания завершения.
//  3. Если буфер заполнен, возвращается ошибка `fmt.Errorf("job buffer is full")`.
//
// Пример использования:
//
//	done, err := service.AddJob(uuid.New(), "# Заголовок\nТекст")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	<-done // Ждём завершения обработки
func (m *PdfService) AddJob(requestID uuid.UUID, inputMarkdown string, user domain.User) (chan struct{}, error) {
	newJob := Job{
		requestID: requestID,
		input:     inputMarkdown,
		user:      user,
		Done:      make(chan struct{}),
	}
	if len(m.jobs) < m.cfg.PdfConfig.JobBufferSize {
		m.jobs <- newJob
		return newJob.Done, nil
	} else {
		return nil, fmt.Errorf("job buffer is full")
	}
}

func (m *PdfService) handleJob(id int) {
	defer m.wg.Done()
	op := "PdfService.handleJob()"
	log := m.logger.With(
		slog.String("op", op),
	)

	log.Info("start pdf job handler")

	for {

		select {
		case <-m.shutdownChannel:
			return
		case job, ok := <-m.jobs:
			if !ok {
				log.Error("failed to send email",
					slog.String("error", "channel is closed"),
					slog.String("requestID", job.requestID.String()),
				)
				return
			}
			var err error
			for retry := range retryCount {
				var e error
				select {
				case <-m.shutdownChannel:
					return
				default:
					e = m.createPdfFromHtml(log, job)
				}
				if e != nil {
					err = e
					log.Error(
						"failed create pdf",
						slog.Int("retry", retry),
						slog.String("error", err.Error()),
						slog.String("requestID", job.requestID.String()),
					)
					time.Sleep(retryDuration)
					continue
				}
				err = e
				break

			}

			if err != nil {
				log.Error(fmt.Sprintf("failed to create pdf after %d retries", retryCount),
					slog.String("error", err.Error()),
					slog.String("requestID", job.requestID.String()),
				)
				return
			}

			m.mailService.AddJob(job.requestID, job.user, "Prof Report")

			close(job.Done)

			log.Info(
				"pdf created",
				slog.Int("id", id),
				slog.String("requestID", job.requestID.String()),
			)

		}
	}
}

func (m *PdfService) createPdf(logger *slog.Logger, job Job) error {

	requestID := job.requestID

	log := logger.With(
		slog.String("op", "PdfService.createPdf()"),
		slog.String("requestID", requestID.String()),
	)

	httpClient := &http.Client{
		Timeout: timeout,
	}

	fullPath := fmt.Sprintf("http://%s:%d", m.cfg.PdfConfig.PdfHost, m.cfg.PdfConfig.PdfPort)

	client, err := gotenberg.NewClient(httpClient, fullPath)
	if err != nil {
		return fmt.Errorf("failed to create gotenberg client: %w", err)
	}

	// indexHTML, err := markdown.FS.ReadFile("/etc/backend-service-profreport/config/template.html")
	// if err != nil {
	// 	return fmt.Errorf("Failed to read template.html: %w", err)
	// }

	// markdownContent, err := markdown.FS.ReadFile("content.md")
	// if err != nil {
	// 	return fmt.Errorf("Failed to read content.md: %w", err)
	// }

	indexHTML, err := os.ReadFile(fmt.Sprintf("%s%s", m.cfg.PdfConfig.HtmlTemplateFilePath, m.cfg.PdfConfig.HtmlTemplateFileName))
	if err != nil {
		return fmt.Errorf("Failed to read template.html: %w", err)
	}

	ctx := context.Background()

	// log.Debug(
	// 	"input data",
	// 	slog.String("indexHTML", string(indexHTML)),
	// 	slog.String("markdownContent", string(markdownContent)),
	// )

	response, err := client.Chromium().
		ConvertMarkdown(ctx, bytes.NewReader(indexHTML)).
		File("content.md", strings.NewReader(job.input)).
		PaperSizeA4().
		Landscape().
		Margins(1, 1, 1, 1).
		OutputFilename(fmt.Sprintf("%s.pdf", requestID.String())).
		Send()

	if err != nil {
		return fmt.Errorf("Failed to convert markdown: %w", err)
	}
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s%s.pdf", m.cfg.PdfConfig.PdfFilePath, requestID.String()))
	if err != nil {
		return fmt.Errorf("Failed to create output file: %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to write PDF: %w", err)
	}

	log.Info(
		"Markdown converted to PDF successfully",
		slog.String("output pdf", fmt.Sprintf("%s.pdf", requestID.String())),
		slog.String("Gotenberg trace: %s", response.GotenbergTrace),
	)

	return nil
}

func (m *PdfService) createPdfFromHtml(logger *slog.Logger, job Job) error {

	requestID := job.requestID

	log := logger.With(
		slog.String("op", "PdfService.createPdfFromHtml()"),
		slog.String("requestID", requestID.String()),
	)

	httpClient := &http.Client{
		Timeout: timeout,
	}

	fullPath := fmt.Sprintf("http://%s:%d", m.cfg.PdfConfig.PdfHost, m.cfg.PdfConfig.PdfPort)

	client, err := gotenberg.NewClient(httpClient, fullPath)
	if err != nil {
		return fmt.Errorf("failed to create gotenberg client: %w", err)
	}

	htmlTmplFullPath := fmt.Sprintf("%s%s", m.cfg.PdfConfig.HtmlTemplateFilePath, m.cfg.PdfConfig.HtmlTemplateFileName)

	user := domain.User{
		Name:  job.user.Name,
		Email: job.user.Email,
	}

	data := ViewData{
		User:           user,
		AiHtmlResponse: job.input,
	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)

	tmpl, _ := template.ParseFiles(htmlTmplFullPath)
	tmpl.Execute(buffer, data)

	ctx := context.Background()

	response, err := client.Chromium().
		ConvertHTML(ctx, buffer).
		PaperSizeA4().
		Landscape().
		Margins(1, 1, 1, 1).
		OutputFilename(fmt.Sprintf("%s.pdf", requestID.String())).
		Send()

	if err != nil {
		return fmt.Errorf("Failed to convert html: %w", err)
	}
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s%s.pdf", m.cfg.PdfConfig.PdfFilePath, requestID.String()))
	if err != nil {
		return fmt.Errorf("Failed to create output file: %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to write PDF: %w", err)
	}

	log.Info(
		"html converted to PDF successfully",
		slog.String("output pdf", fmt.Sprintf("%s.pdf", requestID.String())),
		slog.String("Gotenberg trace: %s", response.GotenbergTrace),
	)

	return nil
}

func (m *PdfService) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("exit Pdf service: %w", ctx.Err())
		default:
			close(m.shutdownChannel)
			close(m.jobs)
			return nil
		}
	}
}
