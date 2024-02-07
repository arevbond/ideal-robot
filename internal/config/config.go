package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"dev"`
	HttpServer `yaml:"http_server"`
	DB         Database
}

type HttpServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Name     string `env:"POSTGRES_DB"`
	Port     int    `env:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	configPath := getEnv("CONFIG_PATH", "")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("configs file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read configs from %s: %v", configPath, err)
	}
	return &cfg

}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
