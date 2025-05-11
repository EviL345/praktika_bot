package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

const configPath = "config.yaml"

type Config struct {
	ChatId   int64  `yaml:"chat_id"`
	BotToken string `yaml:"bot_token"`
	Db       Db     `yaml:"db"`
}

type Db struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
}

func New() *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}
