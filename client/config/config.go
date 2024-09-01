package config

import (
	"github.com/spf13/viper"
	"sync"
	"time"
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
			singleton = &CfgInstance{Config{
				Env: viper.GetString("env"),
				Client: Client{
					Address:      viper.GetString("address"),
					Timeout:      viper.GetDuration("timeout"),
					RetriesCount: viper.GetInt("retriesCount"),
				},
			}}
		})

	return singleton
}
