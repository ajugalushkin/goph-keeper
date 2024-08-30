package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Client struct {
	Address      string        `yaml:"address" env-required:"true"`
	Timeout      time.Duration `yaml:"timeout" env-required:"1h"`
	RetriesCount int           `yaml:"retries_count" env-default:"1"`
}

// Config структура параметров заауска.
type Config struct {
	Env    string `yaml:"env" env-required:"true"`
	Client Client
}

// MustLoad функция заполнения структуры Config, в случае ошибки паникуем.
func MustLoad() *Config {
	configPath := fetchConfig()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(configPath)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path does not exist" + configPath)
	}

	var newConfig Config
	if err := cleanenv.ReadConfig(configPath, &newConfig); err != nil {
		panic("error reading config: " + err.Error())
	}

	return &newConfig
}

// fetchConfig функция для чтения флага config или переменнной окружения CONFIG.
// приоритет flag
func fetchConfig() string {
	var result string

	flag.StringVar(&result, "config", "", "config file path")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG")
	}
	return result
}
