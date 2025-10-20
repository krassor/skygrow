package pdf

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"app/main.go/internal/config"

	"github.com/google/uuid"
	"github.com/nativebpm/gotenberg"
)

const (
	retryCount    int           = 1
	retryDuration time.Duration = 3 * time.Second
	timeout       time.Duration = 30 * time.Second
)

type Job struct {
	requestID uuid.UUID
	inputMd   string
}

// Mailer представляет собой клиент для отправки электронных писем.
type PdfService struct {
	logger          *slog.Logger
	cfg             *config.Config
	jobs            chan Job
	shutdownChannel chan struct{}
	wg              *sync.WaitGroup
}

// NewMailer создает новый экземпляр Mailer.
// Проверяет, что все необходимые параметры присутствуют в конфиге.
func New(logger *slog.Logger, cfg *config.Config) *PdfService {
	return &PdfService{
		logger:          logger,
		cfg:             cfg,
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

// AddJob добавляет новую задачу на отправку письма в очередь.
//
// Параметры:
//   - to: адрес электронной почты получателя (валидируется через net/mail).
//   - subject: тема письма.
//   - body: текстовое содержимое письма.
//
// Возвращает:
//   - nil, если задача успешно добавлена в очередь.
//   - ошибку в следующих случаях:
//     1. Адрес электронной почты `to` имеет некорректный формат.
//     2. Буфер задач (`jobs`) заполнен до максимальной ёмкости, заданной в конфигурации.
func (m *PdfService) AddJob(requestID uuid.UUID, inputMarkdown string) error {

	if len(m.jobs) < m.cfg.PdfConfig.JobBufferSize {
		m.jobs <- Job{
			requestID: requestID,
			inputMd:   inputMarkdown,
		}
		return nil
	} else {
		return fmt.Errorf("job buffer is full")
	}
}

func (m *PdfService) handleJob(id int) {
	defer m.wg.Done()
	op := "PdfService.handleJob()"
	log := m.logger.With(
		slog.String("op", op),
	)
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
					e = m.createPdf(log, job.requestID, job.inputMd)
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

			log.Info(
				"pdf created",
				slog.Int("id", id),
				slog.String("requestID", job.requestID.String()),
			)
		}
	}
}

func (m *PdfService) createPdf(logger *slog.Logger, requestID uuid.UUID, inputMd string) error {

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

	// indexHTML, err := os.ReadFile(fmt.Sprintf("%s%s", m.cfg.PdfConfig.HtmlTemplateFilePath, m.cfg.PdfConfig.HtmlTemplateFileName))
	// if err != nil {
	// 	return fmt.Errorf("Failed to read template.html: %w", err)
	// }

	//markdownContent := []byte(inputMd)

	ctx := context.Background()

	// log.Debug(
	// 	"input data",
	// 	slog.String("indexHTML", string(indexHTML)),
	// 	slog.String("markdownContent", string(markdownContent)),
	// )

	// response, err := client.Chromium().
	// 	ConvertMarkdown(ctx, bytes.NewReader(indexHTML)).
	// 	File("content.md", bytes.NewReader(markdownContent)).
	// 	PaperSizeA4().
	// 	Landscape().
	// 	Margins(1, 1, 1, 1).
	// 	OutputFilename(fmt.Sprintf("%s.pdf", requestID.String())).
	// 	Send()

	htmlReader := strings.NewReader(inputMd)

	response, err := client.Chromium().
		ConvertHTML(ctx, htmlReader).
		PaperSizeA4().
		Landscape().
		Margins(1, 1, 1, 1).
		OutputFilename(fmt.Sprintf("%s.pdf", requestID.String())).
		Send()

	if err != nil {
		return fmt.Errorf("Failed to convert markdown: %w", err)
	}
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s%s.pdf", m.cfg.PdfConfig.HtmlTemplateFilePath, requestID.String()))
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
		slog.String("Gotenberg trace: %s\n", response.GotenbergTrace),
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
