package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type GRPC struct {
	ServerAddress string        `yaml:"server_address" env-required:"true"`
	Timeout       time.Duration `yaml:"timeout" env-default:"1h"`
}

// Config структура параметров заауска.
type Config struct {
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	Env         string        `yaml:"env" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	GRPC        GRPC
	TokenSecret string
}

// MustLoad функция заполнения структуры Config, в случае ошибки паникуем.
func MustLoad() *Config {
	configPath := fetchConfig()
	if configPath == "" {
		panic("config path is empty")
	}

	cfg := MustLoadByPath(configPath)
	cfg.TokenSecret = os.Getenv("TOKEN_SECRET")
	return cfg
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path does not exist" + configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("error reading config: " + err.Error())
	}

	return &config
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
