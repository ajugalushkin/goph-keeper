package models

import "time"

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
}

type LoginParam struct {
	User     string
	Password string
}
