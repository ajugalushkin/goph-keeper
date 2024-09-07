package config

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Storage struct {
	Path string `yaml:"path" env-required:"true"`
}

type GRPC struct {
	Address string        `yaml:"address" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

type Token struct {
	TTL    time.Duration `yaml:"ttl" env-required:"true"`
	Secret string        `yaml:"secret"`
}

type Minio struct {
	Endpoint string `yaml:"endpoint" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	SSL      bool   `yaml:"ssl" env-required:"true"`
	Bucket   string `yaml:"bucket" env-required:"true"`
}

// Config структура параметров заауска.
type Config struct {
	Env     string `yaml:"env" env-required:"true"`
	Storage Storage
	GRPC    GRPC
	Token   Token
	Minio   Minio
}

// MustLoad функция заполнения структуры Config, в случае ошибки паникуем.
func MustLoad() *Config {
	return MustLoadByPath(fetchConfig())
}

func MustLoadByPath(configPath string) *Config {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./server/config")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			slog.Error("Error reading config file: ", slog.String("error", err.Error()))
		}
		slog.Info("Config file not found in ", slog.String("file", configPath))
	} else {
		slog.Info("Using config file: ", slog.String("file", viper.ConfigFileUsed()))
	}

	var newConfig Config
	if err := viper.Unmarshal(&newConfig); err != nil {
		slog.Error("Unable to unmarshal config file: ", slog.String("error", err.Error()))
		panic(err)
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret != "" {
		newConfig.Token.Secret = tokenSecret
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
