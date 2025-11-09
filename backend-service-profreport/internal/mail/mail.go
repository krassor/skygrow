package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/mail"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	"gopkg.in/gomail.v2"

	"app/main.go/internal/config"
	"app/main.go/internal/models/domain"
	"app/main.go/internal/utils/logger/sl"
)

const (
	retryCount    int           = 1
	retryDuration time.Duration = 3 * time.Second
)

type Job struct {
	requestID uuid.UUID
	user      domain.User
	Subject   string
}

// Mailer представляет собой клиент для отправки электронных писем.
type Mailer struct {
	logger          *slog.Logger
	cfg             *config.Config
	jobs            chan Job
	shutdownChannel chan struct{}
	wg              *sync.WaitGroup
}

// NewMailer создает новый экземпляр Mailer.
// Проверяет, что все необходимые параметры присутствуют в конфиге.
func NewMailer(logger *slog.Logger, cfg *config.Config) *Mailer {
	log := logger.With(
		slog.String("op", "mail.NewMailer()"),
	)
	if err := validateConfig(&cfg.MailConfig); err != nil {
		log.Error(
			"failed to create new mailer",
			slog.String("error", err.Error()),
		)
		return nil
	}

	return &Mailer{
		logger:          logger,
		cfg:             cfg,
		jobs:            make(chan Job, cfg.MailConfig.JobBufferSize),
		shutdownChannel: make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}
}

func (s *Mailer) Start() {
	op := "mail.Start()"
	log := s.logger.With(
		slog.String("op", op),
	)
	for i := 0; i < s.cfg.MailConfig.WorkersCount; i++ {
		s.wg.Add(1)
		go s.handleJob(i)
	}
	log.Info(
		"mail service started",
	)

	s.wg.Wait()

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
func (s *Mailer) AddJob(requestID uuid.UUID, user domain.User, subject string) error {
	to := user.Email

	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("Mailer.AddJob(). invalid email address")
	}

	newJob := Job{
		requestID: requestID,
		user:      user,
		Subject:   subject,
	}

	select {
	case <-s.shutdownChannel:
		return fmt.Errorf("service is shutting down")
	default:
		if len(s.jobs) < s.cfg.MailConfig.JobBufferSize {
			s.jobs <- newJob
			return nil
		} else {
			return fmt.Errorf("job buffer is full")
		}
	}
}

// validateConfig проверяет обязательные параметры конфигурации.
func validateConfig(cfg *config.MailConfig) error {
	if cfg.SMTPHost == "" {
		return fmt.Errorf("smtp host is required")
	}
	if cfg.SMTPPort == 0 {
		return fmt.Errorf("smtp port is required")
	}
	if cfg.Username == "" {
		return fmt.Errorf("smtp username is required")
	}
	if cfg.Password == "" {
		return fmt.Errorf("smtp password is required")
	}
	if cfg.FromAddress == "" {
		return fmt.Errorf("from address is required")
	}
	return nil
}

func (s *Mailer) handleJob(id int) {
	defer s.wg.Done()
	op := "mail.handleJob()"
	log := s.logger.With(
		slog.String("op", op),
	)

	log.Info("start mail job handler")

	for {

		select {
		case <-s.shutdownChannel:
			return
		case job, ok := <-s.jobs:

			requestID := job.requestID

			joblog := log.With(
				slog.String("op", op),
				slog.String("requestID", requestID.String()),
			)

			if !ok {
				joblog.Error("failed to send email",
					slog.String("error", "channel is closed"),
					slog.String("to", job.user.Email),
					slog.String("subject", job.Subject),
				)
				return
			}

			var err error

			mailBody := "Здравствуйте, " + job.user.Name + "!\r\n" +
				"По Вашему запросу был сгенерирован отчет\r\n" +
				"Отчет прикреплен к письму во вложении.\r\n" +
				"\r\n\r\nС уважением, команда proffreport."

			mailBody, err = mdToHTML(mailBody)
			if err != nil {
				sl.Err(err)
				break
			}

			for retry := range retryCount {
				var e error
				select {
				case <-s.shutdownChannel:
					return
				default:
					e = s.sendWithGomail(job.requestID, job.user.Email, job.Subject, mailBody)
				}
				if e != nil {
					err = e
					joblog.Error(
						"failed to send email",
						slog.Int("retry", retry),
						slog.String("error", err.Error()),
						slog.String("to", job.user.Email),
						slog.String("subject", job.Subject),
					)
					time.Sleep(retryDuration)
					continue
				}
				err = e
				break

			}

			if err != nil {
				joblog.Error(fmt.Sprintf("failed to send email after %d retries", retryCount),
					slog.String("error", err.Error()),
					slog.String("to", job.user.Email),
					slog.String("subject", job.Subject),
				)
				break
			}

			joblog.Info(
				"mail sended",
				slog.String("to", job.user.Email),
				slog.String("subject", job.Subject),
				slog.Int("id", id),
			)
		}
	}
}

func (s *Mailer) sendWithGomail(requestID uuid.UUID, to string, subject string, body string) error {
	// Создаем временную зону с фиксированным смещением +3 часа (10800 секунд)
	location := time.FixedZone("MSK", 3*3600) // 3 часа = 10800 секунд

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.cfg.MailConfig.FromAddress)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("MIME-Version", "1.0")
	msg.SetHeader("Content-Type", "text/html; charset=\"UTF-8\"")
	msg.SetHeader("Content-Transfer-Encoding", "8bit")
	msg.SetHeader("Date", time.Now().UTC().In(location).Format(time.RFC1123Z))
	msg.SetHeader("Message-ID", fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), s.cfg.MailConfig.FromAddress))
	msg.SetHeader("X-Mailer", "proffreportServiceApp/1.0")
	msg.SetHeader("List-Unsubscribe", fmt.Sprintf("mailto:%s?subject=unsubscribe", s.cfg.MailConfig.FromAddress))

	msg.SetBody("text/html", body)

	msg.Attach(filepath.Clean(fmt.Sprintf("%s%s.pdf", s.cfg.PdfConfig.PdfFilePath, requestID)))

	d := gomail.NewDialer(
		s.cfg.MailConfig.SMTPHost,
		s.cfg.MailConfig.SMTPPort,
		s.cfg.MailConfig.Username,
		s.cfg.MailConfig.Password)

	d.TLSConfig = &tls.Config{
		ServerName: s.cfg.MailConfig.SMTPHost,
	}

	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (s *Mailer) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("exit Mail service: %w", ctx.Err())
		default:
			close(s.shutdownChannel)
			close(s.jobs)
			return nil
		}
	}
}

func mdToHTML(md string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
