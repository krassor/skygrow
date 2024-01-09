package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env-default:"local"`
	StoragePath string           `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServerConfig `yaml:"http_server"`
	DBConfig    DBConfig
	//MigrationsPath string
	//TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

//type GRPCConfig struct {
//	Port    int           `yaml:"port"`
//	Timeout time.Duration `yaml:"timeout"`
//}

type HttpServerConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"USER" env-default:"user"`
	Password string `yaml:"password" env:"PASSWORD"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
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

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
