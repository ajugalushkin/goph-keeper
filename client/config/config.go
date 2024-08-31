package config

import (
	"flag"
	"os"
	"sync"
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
			configPath := fetchConfig()
			if configPath == "" {
				panic("config path is empty")
			}

			var newConfig Config
			if err := cleanenv.ReadConfig(configPath, &newConfig); err != nil {
				panic("error reading config: " + err.Error())
			}

			singleton = &CfgInstance{newConfig}
		})

	return singleton
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
