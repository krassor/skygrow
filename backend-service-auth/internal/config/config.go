package config

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	HttpServer HttpServerConfig `yaml:"http_server" env-required:"true"`
	DBConfig   DBConfig         `yaml:"db" env-required:"true"`
	//MigrationsPath string
	//TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

//type GRPCConfig struct {
//	Port    int           `yaml:"port"`
//	Timeout time.Duration `yaml:"timeout"`
//}

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

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		log.Println("config path is empty. Load default path: \"config/config.yml\"")
		configPath = "config/config.yml"
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res != "" {
		slog.Info("load config path from command line.", "path", res)
		return res
	}
	res = os.Getenv("CONFIG_PATH")
	slog.Info("load config path from env ", "CONFIG_PATH", res)
	return res
}
