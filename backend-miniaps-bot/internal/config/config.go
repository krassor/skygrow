package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

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

	cfg.configPath = configPath
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
		log.Println("load config path from command line.", "path", res)
		return res
	}
	res = os.Getenv("CONFIG_PATH")
	log.Println("load config path from env ", "CONFIG_PATH", res)
	return res
}

func (cfg *Config) Write() error {
	bufWrite, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error config.Write() marshall: %w", err)
	}

	err = os.WriteFile(cfg.configPath, bufWrite, 0775)
	if err != nil {
		return fmt.Errorf("error config.Write() write file: %w", err)
	}
	return nil
}

func (c *AIConfig) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// SetTimeout sets the timeout value
func (c *AIConfig) SetTimeout(timeout time.Duration) {
	c.Timeout = int(timeout.Seconds())
}
