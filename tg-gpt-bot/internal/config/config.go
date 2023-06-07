package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	defaultAdmin           string  = "KrAssor"
	defaultMaxTokens       int     = 512
	defaultTemperature     float32 = 0.5
	defaultN               int     = 1
	defaultSystemRolePromt string  = ""
)

type OpenAIConfig struct {
	SystemRolePromt string  `yaml:"systemRolePromt"`
	MaxTokens       int     `yaml:"maxTokens"`
	Temperature     float32 `yaml:"temperature"`
	N               int     `yaml:"n"`
}

type BotConfig struct {
	Admins []string     `yaml:"admins"`
	OpenAI OpenAIConfig `yaml:"openAI"`
}

type AppConfig struct {
	ConfigFilepath string
	mutex          sync.RWMutex
}

func InitConfig() *AppConfig {
	configFilepath, ok := os.LookupEnv("CONFIG_FILEPATH")
	if !ok {
		log.Warn().Msgf("Cannot find CONFIG_FILEPATH env")
	}

	if configFilepath == "" {
		configFilepath = "config.yml"
	}

	_, err := os.Stat(configFilepath)
	if os.IsNotExist(err) {
		newConfig := BotConfig{
			Admins: []string{defaultAdmin},
			OpenAI: OpenAIConfig{
				SystemRolePromt: defaultSystemRolePromt,
				MaxTokens:       defaultMaxTokens,
				Temperature:     defaultTemperature,
				N:               defaultN,
			},
		}

		_, err := os.Create(configFilepath)
		if err != nil {
			log.Panic().Msgf("Faled to create config file: %s", configFilepath)
		}

		bufWrite, err := yaml.Marshal(newConfig)
		if err != nil {
			log.Panic().Msgf("Error marshalling config struct")
		}
		err = os.WriteFile(configFilepath, bufWrite, 0775)
		if err != nil {
			log.Panic().Msgf("Error writing config file: %s", configFilepath)
		}

		log.Info().Msgf("New config file created: %s", configFilepath)

	} else {
		log.Info().Msgf("Config file already exists")
	}
	return &AppConfig{
		ConfigFilepath: configFilepath,
	}
}

func (config *AppConfig) WriteOpenAIConfig(openAIConfig *OpenAIConfig) error {
	botConfig, err := config.ReadBotConfig()
	if err != nil {
		return fmt.Errorf("Error WriteOpenAIConfig() read botConfig: %w", err)
	}

	botConfig.OpenAI = *openAIConfig

	bufWrite, err := yaml.Marshal(botConfig)
	if err != nil {
		return fmt.Errorf("Error WriteOpenAIConfig() marshall: %w", err)
	}

	config.mutex.Lock()
	defer config.mutex.Unlock()

	err = os.WriteFile(config.ConfigFilepath, bufWrite, 0775)
	if err != nil {
		return fmt.Errorf("Error WriteOpenAIConfig() write file: %w", err)
	}
	return nil
}

func (config *AppConfig) ReadOpenAIConfig() (OpenAIConfig, error) {

	botConfig := BotConfig{}

	config.mutex.RLock()
	defer config.mutex.RUnlock()

	bufRead, err := os.ReadFile(config.ConfigFilepath)
	if err != nil {
		return OpenAIConfig{}, fmt.Errorf("ReadOpenAIConfig() Cannot read config file %s: %w", config.ConfigFilepath, err)
	}

	err = yaml.Unmarshal(bufRead, &botConfig)
	if err != nil {
		return OpenAIConfig{}, fmt.Errorf("ReadOpenAIConfig() Cannot unmarshall config file %s: %w", config.ConfigFilepath, err)
	}

	return botConfig.OpenAI, nil
}

func (config *AppConfig) WriteBotConfig(botConfig *BotConfig) error {
	bufWrite, err := yaml.Marshal(botConfig)
	if err != nil {
		return fmt.Errorf("Error WriteBotConfig() marshall: %w", err)
	}

	config.mutex.Lock()
	defer config.mutex.Unlock()

	err = os.WriteFile(config.ConfigFilepath, bufWrite, 0775)
	if err != nil {
		return fmt.Errorf("Error WriteBotConfig() write file: %w", err)
	}
	return nil
}

func (config *AppConfig) ReadBotConfig() (BotConfig, error) {
	botConfig := BotConfig{}

	config.mutex.RLock()
	defer config.mutex.RUnlock()

	bufRead, err := os.ReadFile(config.ConfigFilepath)
	if err != nil {
		return BotConfig{}, fmt.Errorf("ReadBotConfig() Cannot read config file %s: %w", config.ConfigFilepath, err)
	}

	err = yaml.Unmarshal(bufRead, &botConfig)

	log.Info().Msgf("Read bot config: %v", botConfig)
	log.Info().Msgf("Admins slice: %v", botConfig.Admins)

	if err != nil {
		return BotConfig{}, fmt.Errorf("ReadBotConfig() Cannot unmarshall config file %s: %w", config.ConfigFilepath, err)
	}
	return botConfig, nil
}
