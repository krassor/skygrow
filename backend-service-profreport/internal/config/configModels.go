package config

import "time"

type Config struct {
	Env            string           `yaml:"env" env-default:"local"`
	HttpServer     HttpServerConfig `yaml:"httpServer" env-required:"true"`
	DBConfig       DBConfig         `yaml:"db" env-required:"true"`
	BotConfig      BotConfig        `yaml:"bot" env-required:"true"`
	MailConfig     MailConfig       `yaml:"mail" env-required:"true"`
	PdfConfig      PdfConfig        `yaml:"pdf" env-required:"false"`
	ConfigFilePath string           `yaml:"configFilePath" env:"CONFIG_FILEPATH" env-default:""`
	ConfigFileName string           `yaml:"configFileName" env:"CONFIG_FILENAME" env-default:""`
	configPath     string
}

type HttpServerConfig struct {
	Address string        `yaml:"address" env-required:"true" env-default:"localhost"`
	Port    string        `yaml:"port" env-required:"true" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5"`
	Secret  string        `yaml:"secret" env-required:"true" env-default:"secret"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"DB_USER" env-default:"user"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
}

type AIConfig struct {
	Timeout            int     `yaml:"timeout" env:"AI_TIMEOUT" env-required:"true" env-default:"300"`
	ModelName          string  `yaml:"modelName" env:"AI_MODEL_NAME" env-required:"true"`
	AIApiToken         string  `yaml:"aiapitoken" env:"AI_API_TOKEN" env-required:"true"`
	SystemRolePromt    string  `yaml:"systemRolePromt" env-default:""`
	PromtFilePath      string  `yaml:"promtFilePath" env:"PROMT_FILEPATH" env-required:"true" env-default:""`
	PromtFileName      string  `yaml:"promtFileName" env:"PROMT_FILENAME" env-required:"true" env-default:""`
	AiResponseFilePath string  `yaml:"AiResponseFilePath" env:"AI_RESPONSE_FILEPATH" env-required:"true" env-default:""`
	MaxTokens          int     `yaml:"maxTokens" env-default:"65000"`
	Temperature        float32 `yaml:"temperature" env-default:"0.5"`
	N                  int     `yaml:"n" env-default:"1"`
	JobBufferSize      int     `yaml:"jobBufferSize" env:"AI_BUFFER_SIZE" env-default:"10"`
	WorkersCount       int     `yaml:"workersCount" env:"AI_WORKERS_COUNT" env-default:"1"`
}

type BotConfig struct {
	Admins        []string `yaml:"admins" env-default:"KrAssor"`
	TgbotApiToken string   `yaml:"tgbot_apitoken" env:"TGBOT_APITOKEN" env-required:"true"`
	AI            AIConfig `yaml:"AI"`
}

type MailConfig struct {
	SMTPHost      string `yaml:"smtpHost" env:"SMTP_HOST" env-required:"true" env-default:"smtp.rambler.ru"`
	SMTPPort      int    `yaml:"smtpPort" env:"SMTP_PORT" env-required:"true" env-default:"465"`
	Username      string `yaml:"username" env:"MAIL_USERNAME" env-required:"true" env-default:"proffreport@rambler.ru"`
	Password      string `yaml:"password" env:"MAIL_PASSWORD" env-required:"true" env-default:""`
	FromAddress   string `yaml:"fromAddress" env:"MAIL_FROM_ADDRESS" env-required:"true" env-default:"proffreport@rambler.ru"`
	JobBufferSize int    `yaml:"jobBufferSize" env:"MAIL_JOB_BUFFER_SIZE" env-default:"10"`
	WorkersCount  int    `yaml:"workersCount" env:"MAIL_WORKERS_COUNT" env-default:"3"`
}

type PdfConfig struct {
	PdfHost string `yaml:"pdfHost" env:"PDF_HOST" env-required:"true" env-default:"localhost"`
	PdfPort int    `yaml:"pdfPort" env:"PDF_PORT" env-required:"true" env-default:"3000"`
	// Username      string `yaml:"username" env:"MAIL_USERNAME" env-required:"true" env-default:"proffreport@rambler.ru"`
	// Password      string `yaml:"password" env:"MAIL_PASSWORD" env-required:"true" env-default:""`
	HtmlTemplateFilePath string `yaml:"htmlTemplateFilePath" env:"HTML_TEMPLATE_FILEPATH" env-required:"true" env-default:""`
	HtmlTemplateFileName string `yaml:"htmlTemplateFileName" env:"HTML_TEMPLATE_FILENAME" env-required:"true" env-default:""`
	PdfFilePath          string `yaml:"pdfFilePath" env:"PDF_FILEPATH" env-required:"true" env-default:""`
	JobBufferSize        int    `yaml:"jobBufferSize" env:"PDF_JOB_BUFFER_SIZE" env-default:"10"`
	WorkersCount         int    `yaml:"workersCount" env:"PDF_WORKERS_COUNT" env-default:"3"`
}
