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

// MustLoadByPath is a function that loads and parses a configuration file into a Config struct.
// If a non-empty configPath is provided, it sets the viper configuration file path to the given path.
// Otherwise, it sets the viper configuration file name to "config", type to "yaml", and adds "./server/config" and "." to the search paths.
// It then reads the configuration file and unmarshals it into a new Config struct.
// If an error occurs during the reading or unmarshalling process, it logs the error and panics.
// If the environment variable "TOKEN_SECRET" is set, it updates the Config.Token.Secret field with the value from the environment variable.
// Finally, it returns a pointer to the newly created Config struct.
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

// fetchConfig is a function that reads the configuration file path from command-line flags or environment variables.
// The function prioritizes command-line flags over environment variables.
//
// The function uses the "flag" package to define a command-line flag named "config" with a default value of an empty string.
// The flag is used to specify the path to the configuration file.
//
// If the "config" flag is provided on the command line, the function returns the value of the flag.
// If the "config" flag is not provided, the function checks for the existence of the "CONFIG" environment variable.
// If the "CONFIG" environment variable is set, the function returns the value of the environment variable.
// If the "CONFIG" environment variable is not set, the function returns an empty string.
//
// Parameters:
// None
//
// Return value:
// A string representing the path to the configuration file.
func fetchConfig() string {
	var result string

	flag.StringVar(&result, "config", "", "config file path")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG")
	}
	return result
}
