package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"dev"`
	Port     string         `yaml:"port" env-default:"8080"`
	Database DatabaseConfig `yaml:"database"`
	Nats     NatsConfig     `yaml:"nats"`
}

type DatabaseConfig struct {
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type NatsConfig struct {
	Url        string `yaml:"url"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	ClusterId  string `yaml:"cluster_id"`
	ClientId   string `yaml:"client_id"`
	ProducerId string `yaml:"producer_id"`
	Subject    string `yaml:"subject"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
