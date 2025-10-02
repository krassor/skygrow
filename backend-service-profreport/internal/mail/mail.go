package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/mail"
	"net/smtp"
	"sync"
	"time"

	"gopkg.in/gomail.v2"

	"app/main.go/internal/config"
)

const (
	retryCount    int           = 1
	retryDuration time.Duration = 3 * time.Second
)

type Job struct {
	To      string
	Subject string
	Body    string
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

func (m *Mailer) Start() {
	op := "mail.Start()"
	log := m.logger.With(
		slog.String("op", op),
	)
	for i := 0; i < m.cfg.MailConfig.WorkersCount; i++ {
		m.wg.Add(1)
		go m.handleJob(i)
	}
	log.Info(
		"mail service started",
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
func (m *Mailer) AddJob(to string, subject string, body string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("Mailer.AddJob(). invalid email address")
	}
	if len(m.jobs) < m.cfg.MailConfig.JobBufferSize {
		m.jobs <- Job{
			To:      to,
			Subject: subject,
			Body:    body,
		}
		return nil
	} else {
		return fmt.Errorf("job buffer is full")
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

// Send отправляет электронное письмо на указанный адрес.
func (m *Mailer) handleJob(id int) {
	defer m.wg.Done()
	op := "mail.handleJob()"
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
					slog.String("to", job.To),
					slog.String("subject", job.Subject))
				return
			}
			var err error
			for retry := range retryCount {
				var e error
				select {
				case <-m.shutdownChannel:
					return
				default:
					e = m.sendWithGomail(job.To, job.Subject, job.Body)
				}
				if e != nil {
					err = e
					log.Error(
						"failed to send email",
						slog.Int("retry", retry),
						slog.String("error", err.Error()),
						slog.String("to", job.To),
						slog.String("subject", job.Subject),
					)
					time.Sleep(retryDuration)
					continue
				}
				err = e
				break

			}

			if err != nil {
				log.Error(fmt.Sprintf("failed to send email after %d retries", retryCount),
					slog.String("error", err.Error()),
					slog.String("to", job.To),
					slog.String("subject", job.Subject))
				return
			}

			log.Info(
				"mail sended",
				slog.String("to", job.To),
				slog.String("subject", job.Subject),
				slog.Int("id", id),
			)
		}
	}
}

// sendv3 sends an email using SMTP with TLS.
func (m *Mailer) send(to string, subject string, body string) error {
	fromAddress := m.cfg.MailConfig.FromAddress
	username := m.cfg.MailConfig.Username
	password := m.cfg.MailConfig.Password
	smtpHost := m.cfg.MailConfig.SMTPHost
	smtpPort := m.cfg.MailConfig.SMTPPort
	smtpAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	// Формирование заголовков письма
	header := fmt.Sprintf("From: %s\r\n", fromAddress) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject)

	// Формирование тела письма
	message := header + "\r\n" + body + "\r\n"

	// Создание TLS-конфигурации
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	// Установка TLS-соединения
	conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	// Создание SMTP-клиента
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer c.Quit()

	// // Start TLS
	// if err := c.StartTLS(tlsConfig); err != nil {
	// 	return fmt.Errorf("failed to start TLS. c.StartTLS: %w", err)
	// }

	// Auth
	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set the sender and recipient first
	if err := c.Mail(fromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %w", err)
	}
	defer wc.Close()

	_, err = fmt.Fprint(wc, message)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	return nil
}

func (m *Mailer) sendWithGomail(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.MailConfig.FromAddress)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Prof report")
	msg.SetHeader("MIME-Version", "1.0")
	msg.SetHeader("Content-Type", "text/html; charset=\"UTF-8\"")
	msg.SetHeader("Content-Transfer-Encoding", "8bit")
	msg.SetHeader("Date", time.Now().Format(time.RFC1123Z))
	msg.SetHeader("Message-ID", fmt.Sprintf("<%d.%s>", time.Now().UnixNano(), m.cfg.MailConfig.FromAddress))
	msg.SetHeader("X-Mailer", "proffreportServiceApp/1.0")

	msg.SetBody("text/html", body)

	d := gomail.NewDialer(
		m.cfg.MailConfig.SMTPHost,
		m.cfg.MailConfig.SMTPPort,
		m.cfg.MailConfig.Username,
		m.cfg.MailConfig.Password)

	d.TLSConfig = &tls.Config{
		ServerName: m.cfg.MailConfig.SMTPHost,
	}

	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (m *Mailer) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("exit Mail service: %w", ctx.Err())
		default:
			close(m.shutdownChannel)
			close(m.jobs)
			return nil
		}
	}
}
