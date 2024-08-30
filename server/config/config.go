package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Token struct {
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	TokenSecret string        `yaml:"token_secret"`
}

type GRPC struct {
	ServerAddress string        `yaml:"server_address" env-required:"true"`
	Timeout       time.Duration `yaml:"timeout" env-default:"1h"`
}

// Config структура параметров заауска.
type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	GRPC        GRPC
	Token       Token
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

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret != "" {
		newConfig.Token.TokenSecret = tokenSecret
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
