package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/ajugalushkin/goph-keeper/cmd"
)

// Config структура параметров заауска.
type Config struct {
	ServerAddress string        `env:"SERVER_ADDRESS" env-required:"true"`
	TokenTTL      time.Duration `env:"TOKEN_TTL" env-required:"true"`
	Env           string        `env:"ENV" env-required:"true"`
}

// init функция инициализации начальных значений для параметров запуска.
func init() {
	err := godotenv.Load("/docker/.env")
	if err != nil {
		//log.Debug("Error loading .env file", "error", err)
	}
	viper.SetDefault("Server_Address", "localhost:8080")
	viper.SetDefault("Token_TTL", "1h")
	viper.SetDefault("Env", "")
}

// bindToEnv функция для маппинга полей из ENV с полями структуры.
func bindToEnv() {
	_ = viper.BindEnv("Server_Address")
	_ = viper.BindEnv("Token_TTL")
	_ = viper.BindEnv("Env")
}

func MustLoad() *Config {
	bindToEnv()

	err := cmd.Execute()
	if err != nil {
		//fmt.Println(err)
	}

	result := Config{
		ServerAddress: viper.GetString("Server_Address"),
		TokenTTL:      viper.GetDuration("Token_TTL"),
		Env:           viper.GetString("Env"),
	}

	if result.ServerAddress == "" {
		panic("Server_Address is empty")
	}

	if result.TokenTTL == 0 {
		panic("Token_TTL is empty")
	}

	if result.Env == "" {
		panic("Env is empty")
	}

	return &result
}
