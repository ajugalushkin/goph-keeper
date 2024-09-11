package config

import (
	"log/slog"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Client struct {
	Address string        `yaml:"address" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
	Retries int           `yaml:"retries" env-required:"true"`
}

// Config структура параметров заауска.
type Config struct {
	Env    string `yaml:"env" env-required:"true"`
	Client Client `yaml:"client" env-required:"true"`
}

type CfgInstance struct {
	Config Config
}

var (
	singleton *CfgInstance
	once      sync.Once
)

func GetInstance() *CfgInstance {
	once.Do(
		func() {
			var config Config
			if err := viper.Unmarshal(&config); err != nil {
				slog.Error("Unable to unmarshal config file: ", slog.String("error", err.Error()))
				panic(err)
			}

			singleton = &CfgInstance{config}
		})

	return singleton
}
